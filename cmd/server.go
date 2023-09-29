package main

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"processing-orchestrator/pkg/cooker"
	"processing-orchestrator/pkg/encoder"
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	"processing-orchestrator/pkg/uploader"
	pb "processing-orchestrator/proto"
	progress_reporter "processing-orchestrator/services/progress-reporter"
	"processing-orchestrator/services/record-processor"
	"strconv"
)

const (
	DEFAULT_PORT      = 55556
	DEFAULT_DAPR_PORT = 50001
	// Dapr services app ids
	// TODO :: Move these to env vars
	DEFAULT_UPLOADER_ID           = "video-store"
	DEFAULT_PUB_COMPONENT         = "message-queue"
	DEFAULT_SUB_COMPONENT         = "pubsub"
	DEFAULT_STATE_STORE_COMPONENT = "statestore"
)

type server struct {
	pb.UnimplementedProcessorServer
	processor   *record_processor.RecordProcessor
	progressRep *progress_reporter.ProgressReporter
}

func (s *server) Watch(req *pb.WatchRequest, stream pb.Processor_WatchServer) error {
	com := &progress_reporter.BidirectionalCom{
		Data:   make(chan *pb.ProcessingStatus, 20),
		Canary: make(chan bool, 1),
	}
	defer close(com.Canary)
	defer close(com.Data)
	s.progressRep.Register(req.Id, com)

	for {
		progress := <-com.Data
		if err := stream.Send(progress); err != nil {
			// Client closed the stream
			if err == io.EOF {
				break
			}
			log.Printf("Error sending status: %v", err)
			return err
		}
		// Server will close the stream
		if progress.Done {
			break
		}
	}
	return nil
}

func (s *server) Start(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	jobId, err := s.processor.Add(req.DiscordAudioKeys, req.BackgroundAudioKey)
	if err != nil {
		return nil, err
	}
	s.progressRep.AddInfo(jobId, req.DiscordAudioKeys)
	slog.Info(fmt.Sprintf("[Server] :: Starting a new process with id '%s' and params %+v", jobId, req))
	return &pb.ProcessResponse{Id: jobId}, nil
}

func main() {
	pEnv := parseEnv()
	slog.Info("[Main] :: Dapr port is " + strconv.Itoa(pEnv.daprGrpcPort))

	// Strat the gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", pEnv.serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	daprServer := daprd.NewServiceWithGrpcServer(lis, s)
	processor, reporter, err := DI(daprServer, pEnv)
	if err != nil {
		log.Fatalf("failed to initialize event controller: %v", err)
	}
	err = processor.Init()
	go reporter.Start()

	if err != nil {
		log.Fatalf("[Main] :: Processor failed to init: %v", err)
	}
	pb.RegisterProcessorServer(s, &server{processor: processor, progressRep: reporter})
	slog.Info(fmt.Sprintf("[Main] :: Starting gRPC server at %v", lis.Addr()))
	if err := daprServer.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}

type env struct {
	// Port to connect to Dapr sidecar
	daprGrpcPort int
	// Port the app is listening on
	serverPort int
	// Dapr components ids
	daprCpnUploader string
	daprCpnPub      string
	daprCpnSub      string
	daprCpnState    string
}

func parseEnv() *env {
	pEnv := env{
		serverPort:      DEFAULT_PORT,
		daprGrpcPort:    DEFAULT_DAPR_PORT,
		daprCpnUploader: DEFAULT_UPLOADER_ID,
		daprCpnPub:      DEFAULT_PUB_COMPONENT,
		daprCpnSub:      DEFAULT_SUB_COMPONENT,
		daprCpnState:    DEFAULT_STATE_STORE_COMPONENT,
	}

	if envPort, err := strconv.ParseInt(os.Getenv("DAPR_GRPC_PORT"), 10, 32); err == nil && envPort != 0 {
		pEnv.daprGrpcPort = int(envPort)
	}
	if envPort, err := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 32); err == nil && envPort != 0 {
		pEnv.serverPort = int(envPort)
	}
	if id, isDefined := os.LookupEnv("UPLOADER_NAME"); isDefined && id != "" {
		pEnv.daprCpnUploader = id
	}
	if id, isDefined := os.LookupEnv("PUBSUB_NAME"); isDefined && id != "" {
		pEnv.daprCpnPub = id
	}
	if id, isDefined := os.LookupEnv("SUB_NAME"); isDefined && id != "" {
		pEnv.daprCpnSub = id
	}
	if id, isDefined := os.LookupEnv("STORE_NAME"); isDefined && id != "" {
		pEnv.daprCpnState = id
	}
	return &pEnv
}

func DI(subServer common.Service, env *env) (*record_processor.RecordProcessor, *progress_reporter.ProgressReporter, error) {
	daprClient, err := makeDaprClient(env.daprGrpcPort, 16)
	if err != nil {
		return nil, nil, err
	}
	progressCh := make(chan processing_common.Watchable, 100)
	cook := cooker.NewCooker(daprClient, env.daprCpnPub, env.daprCpnSub, progressCh)
	err = cook.SubscribeTo(subServer)
	if err != nil {
		return nil, nil, err
	}
	encode := encoder.NewEncoder(daprClient, env.daprCpnPub, env.daprCpnSub, progressCh)
	err = encode.SubscribeTo(subServer)
	if err != nil {
		return nil, nil, err
	}
	upload := uploader.NewUploader(daprClient, env.daprCpnUploader, env.daprCpnSub, progressCh)
	err = upload.SubscribeTo(subServer)
	if err != nil {
		return nil, nil, err
	}
	store := job_store.NewJobStore(daprClient, env.daprCpnState)
	reporter := progress_reporter.NewProgressReporter(progressCh)
	return record_processor.NewRecordProcessor(cook, encode, upload, store), reporter, nil
}

func makeDaprClient(port, maxRequestSizeMB int) (client.Client, error) {
	var opts []grpc.CallOption
	opts = append(opts, grpc.MaxCallRecvMsgSize(maxRequestSizeMB*1024*1024))
	conn, err := grpc.Dial(net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)),
		grpc.WithDefaultCallOptions(opts...), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return client.NewClientWithConnection(conn), nil
}

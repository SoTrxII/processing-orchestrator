package main

import (
	"context"
	"flag"
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
	DEFAULT_DAPR_PORT = 50001
	// Dapr services app ids
	// TODO :: Move these to env vars
	DEFAULT_UPLOADER_ID           = "video-store"
	DEFAULT_PUBSUB_COMPONENT      = "pubsub"
	DEFAULT_STATE_STORE_COMPONENT = "statestore"
)

var (
	port = flag.Int("port", 55556, "The server port")
)

type server struct {
	pb.UnimplementedRecordServiceServer
	processor   *record_processor.RecordProcessor
	progressRep *progress_reporter.ProgressReporter
}

func (s *server) Watch(req *pb.WatchRequest, stream pb.RecordService_WatchServer) error {
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
	return &pb.ProcessResponse{Id: jobId}, nil
}

func main() {
	daprPort := DEFAULT_DAPR_PORT
	if envPort, err := strconv.ParseInt(os.Getenv("DAPR_GRPC_PORT"), 10, 32); err == nil && envPort != 0 {
		daprPort = int(envPort)
	}
	slog.Info("[Main] :: Dapr port is " + strconv.Itoa(daprPort))

	// Strat the gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	daprServer := daprd.NewServiceWithGrpcServer(lis, s)
	processor, err := DI(daprServer, daprPort)
	if err != nil {
		log.Fatalf("failed to initialize event controller: %v", err)
	}
	err = processor.Init()
	if err != nil {
		log.Fatalf("[Main] :: Processor failed to init: %v", err)
	}
	pb.RegisterRecordServiceServer(s, &server{processor: processor})
	slog.Info(fmt.Sprintf("[Main] :: Starting gRPC server at %v", lis.Addr()))
	if err := daprServer.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}

func DI(subServer common.Service, daprPort int) (*record_processor.RecordProcessor, error) {
	daprClient, err := makeDaprClient(daprPort, 16)
	if err != nil {
		return nil, err
	}
	progressCh := make(chan processing_common.Watchable, 100)
	cook := cooker.NewCooker(daprClient, subServer, DEFAULT_PUBSUB_COMPONENT, progressCh)
	encode := encoder.NewEncoder(daprClient, subServer, DEFAULT_PUBSUB_COMPONENT, progressCh)
	upload := uploader.NewUploader(daprClient, DEFAULT_UPLOADER_ID)
	store := job_store.NewJobStore(daprClient, DEFAULT_STATE_STORE_COMPONENT)
	return record_processor.NewRecordProcessor(cook, encode, upload, store), nil
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

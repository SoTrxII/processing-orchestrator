package thumb_generator

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	pb "processing-orchestrator/pkg/thumb-generator/proto"
	"time"
)

// A thumbnail takes ~3s (2026-10-18), most of it downloading the background.
// The upload waits for it, so without a bound a generator that never answers
// holds a finished video forever: on 2026-10-18 one sat after rendering the
// image, 12 min until its pod was restarted. Failing it only costs the
// thumbnail, which the upload already treats as optional
const generationTimeout = 2 * time.Minute

type ThumbGenerator struct {
	client          pb.ThumbnailClient
	ctx             context.Context
	invokeComponent string
	timeout         time.Duration
}

func NewThumbGenerator(daprAddress, appId string) *ThumbGenerator {
	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, daprAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	methodCtx := metadata.AppendToOutgoingContext(context.Background(), "dapr-app-id", appId)
	methodCtx = metadata.AppendToOutgoingContext(methodCtx, "dapr-stream", "true")

	return &ThumbGenerator{
		client:  pb.NewThumbnailClient(conn),
		ctx:     methodCtx,
		timeout: generationTimeout,
	}
}

func (t *ThumbGenerator) GenerateThumbnail(req *pb.ThumbnailRequest) (string, error) {
	ctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	res, err := t.client.CreateThumbnail(ctx, req)
	if err != nil {
		return "", err
	}
	return res.ThumbnailKey, nil
}

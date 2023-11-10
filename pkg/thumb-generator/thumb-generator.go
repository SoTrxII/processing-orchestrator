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

type ThumbGenerator struct {
	client          pb.ThumbnailClient
	ctx             context.Context
	invokeComponent string
}

func NewThumbGenerator(daprAddress, appId string) *ThumbGenerator {
	dialCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(dialCtx, daprAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	methodCtx := metadata.AppendToOutgoingContext(context.Background(), "dapr-app-id", appId)
	methodCtx = metadata.AppendToOutgoingContext(methodCtx, "dapr-stream", "true")

	return &ThumbGenerator{
		client: pb.NewThumbnailClient(conn),
		ctx:    methodCtx,
	}
}

func (t *ThumbGenerator) GenerateThumbnail(req *pb.ThumbnailRequest) (string, error) {
	res, err := t.client.CreateThumbnail(t.ctx, req)
	if err != nil {
		return "", err
	}
	return res.ThumbnailKey, nil
}

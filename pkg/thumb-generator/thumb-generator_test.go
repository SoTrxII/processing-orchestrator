package thumb_generator

import (
	"context"
	"google.golang.org/grpc"
	pb "processing-orchestrator/pkg/thumb-generator/proto"
	"testing"
	"time"
)

// Answers only when the caller gives up, like the generator of 2026-10-18
type silentClient struct{ pb.ThumbnailClient }

func (silentClient) CreateThumbnail(ctx context.Context, _ *pb.ThumbnailRequest, _ ...grpc.CallOption) (*pb.ThumbnailResponse, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestGenerateThumbnail_ASilentGeneratorTimesOut(t *testing.T) {
	gen := &ThumbGenerator{client: silentClient{}, ctx: context.Background(), timeout: 10 * time.Millisecond}
	done := make(chan error, 1)
	go func() { _, err := gen.GenerateThumbnail(&pb.ThumbnailRequest{}); done <- err }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("a generator that never answers must fail the call")
		}
	case <-time.After(time.Second):
		t.Fatal("the call waited on a silent generator")
	}
}

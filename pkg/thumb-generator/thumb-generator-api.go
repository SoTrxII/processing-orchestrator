package thumb_generator

import (
	"context"

	pb "processing-orchestrator/pkg/thumb-generator/proto"
)

type ThumbnailRequest = pb.ThumbnailRequest
type ThumbGenService interface {
	GenerateThumbnail(ctx context.Context, req *ThumbnailRequest) (string, error)
}

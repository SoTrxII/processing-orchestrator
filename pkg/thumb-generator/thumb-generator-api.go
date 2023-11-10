package thumb_generator

import pb "processing-orchestrator/pkg/thumb-generator/proto"

type ThumbnailRequest = pb.ThumbnailRequest
type ThumbGenService interface {
	GenerateThumbnail(req *ThumbnailRequest) (string, error)
}

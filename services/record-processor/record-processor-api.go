package record_processor

import thumb_generator "processing-orchestrator/pkg/thumb-generator"

// Addons, if any, are services that are not part of the core processing pipeline
type Addons struct {
	ThumbGen thumb_generator.ThumbGenService
}

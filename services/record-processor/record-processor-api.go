package record_processor

import (
	"processing-orchestrator/pkg/summarizer"
	thumb_generator "processing-orchestrator/pkg/thumb-generator"
)

// Addons, if any, are services that are not part of the core processing pipeline
type Addons struct {
	ThumbGen thumb_generator.ThumbGenService
	// Summarizer turns the recording into a narrative summary. It runs
	// entirely on its own clock, well after this job is done.
	Summarizer summarizer.SummarizingService
}

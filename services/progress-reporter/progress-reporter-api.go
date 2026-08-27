package progress_reporter

import (
	processing_common "processing-orchestrator/pkg/processing-common"
	pb "processing-orchestrator/proto"
)

var steps = []string{
	processing_common.StepCooking.ToString(),
	processing_common.StepEncoding.ToString(),
	processing_common.StepUploading.ToString(),
}

// BidirectionalCom is how the reporter talks to one client watching a job.
// Only Done is ever closed, and only by the client, to say it stopped
// watching : closing Data on the client side would leave the reporter sending
// on a closed channel, taking the whole process down with it
type BidirectionalCom struct {
	// Progress of the job, written by the reporter only
	Data chan *pb.ProcessingStatus
	// Closed by the client once it stops watching
	Done chan struct{}
}

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

type ProgressReporter struct {
	// Maps id to audioskeys
	jobInfos map[string][]string
	// maps id to watch server
	jobWatchers map[string][]*BidirectionalCom
}

type BidirectionalCom struct {
	Data   chan *pb.ProcessingStatus
	Canary chan bool
}

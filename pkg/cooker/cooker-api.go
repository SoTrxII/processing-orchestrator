package cooker

import (
	"fmt"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type CookingService interface {
	Cook(jobId string, recordIds []string) ([]string, error)
}

const (
	// Topic watched by the cooking service
	p_Start = "cookings"
	// Topic the cooking service is publishing into
	s_Info = "cooking-state"
)

type CookerOpt struct {
	// File extension of the cooked record.
	// This is normally supplied with each done event
	// Default is '.ogg'
	// The '.' MUST be supplied
	Extension string
}

// CookingJob Arguments that have to be passed to the cooking service
type CookingJob struct {
	// UUID, can be anything
	JobId string `json:"jobId"`
	// Keys of the discord recording.
	// As a single recording can consist of multiple
	// files (in case of DR), there can be multiples keys
	Ids []string `json:"ids"`
}

type CookingData struct {
	// Size of the transcoded record (bytes)
	// Only supplied if cooking state is Progress
	TotalBytes int64 `json:"totalBytes"`
	// Error message if any
	// Only supplied if cooking state is Error
	Message string `json:"message"`
	// File Extension of the record (with the ".")
	// Only supplied if cooking state is Done
	Extension string `json:"Extension"`
}

type CookingEvent struct {
	processing_common.ServiceEvent
	// Current pandora record being processed
	RecordId string      `json:"recordId"`
	Data     CookingData `json:"data"`
}

// ToProgress() Return a user-readable representation of what's going on
func (ce *CookingEvent) ToProgress() *processing_common.ServiceProgress {

	pg := processing_common.ServiceProgress{
		JobId:       ce.JobId,
		Step:        processing_common.StepCooking,
		CurrentItem: ce.RecordId,
		Error:       nil,
		Progress:    "",
	}
	switch ce.State {
	case processing_common.Done:
		pg.Progress = "Done"
	case processing_common.InProgress:
		pg.Progress = fmt.Sprintf("%d bytes processed so far", ce.Data.TotalBytes)
	case processing_common.Error:
		pg.Error = fmt.Errorf(ce.Data.Message)
	default:
		pg.Error = fmt.Errorf("unsupported state %d", ce.State)
	}
	return &pg
}

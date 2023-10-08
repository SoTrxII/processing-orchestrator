package encoder

import (
	"fmt"
	"math"
	processing_common "processing-orchestrator/pkg/processing-common"
	"strings"
	"time"
)

type EncodingService interface {
	Encode(jobId string, audioKeys []string, bgAudioKey string) (string, error)
}

const (
	// Topic watched by the encoding service
	p_Start = "encodings"
	// Topic the encoding service is publishing into
	s_Info = "encoding-state"
)

type EncoderOpt struct {
	// File extension of the encoded record.
	// This is normally supplied with each done event
	// Default is '.mp4'
	// The '.' MUST be supplied
	Extension string
}
type EncodeJobOpt struct {
	/** Whether to delete used assets (videoKey, audioKeys and ImageKey)
	 * from the remote object storage. Default is false */
	deleteAssetsFromObjStore bool
}
type EncodeJob struct {
	/** UUID of the processing job */
	JobId string `json:"jobId"`
	/** Storage backend retrieval keys for all videos tracks */
	VideoKey string `json:"videoKey" omitempty:"true"`
	/** Storage backend retrieval keys for all audio tracks */
	AudiosKeys []string `json:"audiosKeys"`
	// Storage backend keys for all side audio track part
	BackgroundAudioKey string `json:"backgroundAudioKey" omitempty:"true"`
	/** Storage backend retrieval keys for the image track */
	ImageKey string `json:"imageKey" omitempty:"true"`
	/** Options for the encoding job */
	Opt EncodeJobOpt `json:"options"`
}

type EncodingData struct {
	// Number of total frames processed
	Frames int64 `json:"frames"`
	// Number of frames processed each second
	Fps int `json:"fps"`
	// Quality target. Usually between 20 and 30
	Quality float32 `json:"quality"`
	// Estimated size of the converted file (kb)
	Size int64 `json:"size"`
	// Total processed time
	Time time.Duration `json:"time"`
	// Target bitrate
	Bitrate string `json:"bitrate"`
	// Encoding speed. A "2" means 1 second of encoding would be a 2 seconds playback
	Speed float32 `json:"speed"`
	// The duration of the output file
	TargetDuration time.Duration `json:"totalDuration"`
	// Error message
	// Only set if the state is "error"
	Message string `json:"message"`
}

// Returns the percentage of the encoding progress
func (ed *EncodingData) GetProgressPercentage() (float64, error) {
	// Preventing a div by 0 zero
	if ed.TargetDuration == 0 {
		return 0, fmt.Errorf("target duration is 0")
	}

	return math.Round((ed.Time.Seconds() / ed.TargetDuration.Seconds()) * 100), nil
}

type EncodingEvent struct {
	processing_common.ServiceEvent
	Data EncodingData `json:"data"`
}

func (e *EncodingEvent) ToProgress() *processing_common.ServiceProgress {
	pg := processing_common.ServiceProgress{
		JobId:       e.JobId,
		Step:        processing_common.StepEncoding,
		CurrentItem: "video",
		Error:       nil,
		Progress:    "",
	}
	switch e.State {
	case processing_common.Done:
		pg.Progress = "Done"
	case processing_common.Error:
		pg.Error = fmt.Errorf(e.Data.Message)
	case processing_common.InProgress:
		ss := strings.Builder{}
		ss.WriteString(fmt.Sprintf("Processed %s", e.Data.Time))
		if e.Data.TargetDuration != 0 {
			percentage := math.Min(e.Data.Time.Seconds()/e.Data.TargetDuration.Seconds()*100, 100)
			ss.WriteString(fmt.Sprintf(" on total %s (%.2f%%)", e.Data.TargetDuration, percentage))
		}
		pg.Progress = ss.String()
	default:
		pg.Error = fmt.Errorf("Unsupported state %d", e.State)
	}
	return &pg
}

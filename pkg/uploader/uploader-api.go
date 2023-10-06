package uploader

import (
	"fmt"
	processing_common "processing-orchestrator/pkg/processing-common"
	"time"
)

type UploadingService interface {
	Upload(jobId, storageKey string, opt *processing_common.VideoOpt) (*Video, error)
}

const (
	// Topic the upload service is publishing into
	s_Info = "upload-state"
)

type UploadJob struct {
	processing_common.VideoOpt
	/** Key to retrieve the video assets on the remote object storage */
	StorageKey string `json:"storageKey"`
	/** Job UUID */
	JobId string `json:"jobId"`
}

type Video struct {
	Id string `json:"id"`
	// Video display name
	Title string `json:"title" validate:"updatable"`
	// Video description
	Description string `json:"description" validate:"updatable"`
	// Creation date
	CreatedAt time.Time `json:"createdAt"`
	// Video duration in seconds
	Duration int64 `json:"duration"`
	// public/private/unlisted
	Visibility processing_common.Visibility `json:"visibility" validate:"updatable"`
	// Playlist thumbnail
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	// Url prefix necessary to watch the video. ie https://www.youtube.com/watch?v= for Youtube
	WatchPrefix string `json:"watchPrefix"`
}

type UploadData struct {
	/** Number of bytes uploaded so far */
	Current int64 `json:"current"`
	/** Size of the file to upload (bytes)*/
	Total int64 `json:"total"`
	/** Error message if any */
	Message string `json:"message"`
}
type UploadEvent struct {
	processing_common.ServiceEvent
	Data UploadData `json:"data"`
}

func (e *UploadEvent) ToProgress() *processing_common.ServiceProgress {
	pg := processing_common.ServiceProgress{
		JobId:       e.JobId,
		Step:        processing_common.StepUploading,
		CurrentItem: "video",
		Error:       nil,
		Progress:    "",
	}
	switch e.State {
	case processing_common.Done:
		pg.Progress = "Done"
	case processing_common.Error:
		// TODO, Propagate error from upstream (not sent yet)
		pg.Error = fmt.Errorf("Unkown error")
	case processing_common.InProgress:
		// TODO : Upsert ETA
		pg.Progress = fmt.Sprintf("Uploaded %d bytes on total %d", e.Data.Current, e.Data.Total)
	default:
		pg.Error = fmt.Errorf("unsupported state %d", e.State)
	}
	return &pg
}

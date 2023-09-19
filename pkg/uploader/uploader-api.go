package uploader

import "time"

type UploadingService interface {
	Upload(jobId, storageKey string, opt *VideoOpt) (*Video, error)
}

type Visibility string

const (
	Public   Visibility = "public"
	Private  Visibility = "private"
	Unlisted Visibility = "unlisted"
)

type VideoOpt struct {
	// Short text describing the content of the item
	// Youtube actually limits to 5000 bytes, which *isn't* 5000 characters
	// https://developers.google.com/youtube/v3/docs/videos#properties
	Description string `json:"description" binding:"max=1000"`
	// Title of the item
	// The max character limitation is currently taken from the Yt docs
	// https://developers.google.com/youtube/v3/docs/videos#properties
	// This may change if another provider is requiring less than 100 characters
	Title string `json:"title" binding:"required,max=100"`
	// Visibility of the item
	Visibility Visibility `json:"visibility" binding:"required"`
}

type UploadJob struct {
	VideoOpt
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
	Visibility Visibility `json:"visibility" validate:"updatable"`
	// Playlist thumbnail
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	// Url prefix necessary to watch the video. ie https://www.youtube.com/watch?v= for Youtube
	WatchPrefix string `json:"watchPrefix"`
}

package processing_common

type State int8

const (
	InProgress State = iota
	Done
	Error
)

type ServiceEvent struct {
	JobId string `json:"jobId"`
	State State  `json:"state"`
}

type Watchable interface {
	ToProgress() *ServiceProgress
}

type Step int8

const (
	StepCooking Step = iota
	StepEncoding
	StepUploading
	StepDone
)

func (s Step) ToString() string {
	switch s {
	case StepCooking:
		return "Cooking"
	case StepEncoding:
		return "Encoding"
	case StepUploading:
		return "Uploading"
	case StepDone:
		return "Done"
	default:
		return "Unknown"
	}
}

type ServiceProgress struct {
	JobId       string
	Step        Step
	CurrentItem string
	Error       error
	Progress    string
	Link        string
}

type DoneData struct {
	Link string `json:"link"`
}

type DoneEvent struct {
	ServiceEvent
	Data DoneData `json:"data"`
}

func (e *DoneEvent) ToProgress() *ServiceProgress {
	pg := ServiceProgress{
		JobId:       e.JobId,
		Step:        StepDone,
		CurrentItem: "",
		Error:       nil,
		Progress:    "",
		Link:        e.Data.Link,
	}
	return &pg
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
	// Optional playlist to put the newly uploaded video intp
	PlaylistId string `json:"playlistId" omitempty:"true"`
	// Mandatory playlist title. If playlistId is not provided, a new playlist will be created with this title
	PlaylistTitle string `json:"playlistTitle" binding:"required"`
	// Optional thumbnail
	Thumbnail ThumbnailOpt `json:"thumbnail" omitempty:"true"`
}

type ThumbnailOpt struct {
	// Main title, at the center of the thumbnail
	Title string `json:"title" omitempty:"true"`
	// Subtitle, at the bottom of the thumbnail
	SubTitle string `json:"subTitle" omitempty:"true"`
	// Optional numbering
	Number int `json:"number" omitempty:"true"`
	// Background image
	BgUrl string `json:"bgUrl" omitempty:"true"`
}

type UserInput struct {
	Vid VideoOpt
}

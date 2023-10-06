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

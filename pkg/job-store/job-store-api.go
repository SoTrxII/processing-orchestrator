package job_store

import processing_common "processing-orchestrator/pkg/processing-common"

type JobStore interface {
	Upsert(job *JobState) error
	Get(id string) (*JobState, error)
	GetAll(reload bool) ([]*JobState, error)
	Delete(id string) error
}

type JobState struct {
	Id string
	// Cooking, Encoding, uploading, Done
	Step processing_common.Step
	// Audio as recorded by Pandora
	RawAudioKeys []string
	// Audio as cooked by the Cooker
	CookedAudioKeys []string
	// Background audio, optional
	BackgroundAudioKey string
	// Video as encoded by the Encoder
	VideoKey string
	// User provided data
	UserInput processing_common.UserInput
	// Video Link as uploaded by the Uploader
	VideoLink string
	// Link to the created playlist if any
	CreatedPlaylistLink string
	// W3C trace context of the call that started the job, carried by every
	// call made for it (utils.TraceContext). Empty for jobs stored before
	TraceParent string
}

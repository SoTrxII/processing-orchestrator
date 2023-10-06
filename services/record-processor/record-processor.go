package record_processor

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"processing-orchestrator/pkg/cooker"
	"processing-orchestrator/pkg/encoder"
	job_store "processing-orchestrator/pkg/job-store"
	"processing-orchestrator/pkg/processing-common"
	"processing-orchestrator/pkg/uploader"
)

type RecordProcessor struct {
	cooker     cooker.CookingService
	encoder    encoder.EncodingService
	uploader   uploader.UploadingService
	progressCh chan processing_common.Watchable
	store      job_store.JobStore
}

func NewRecordProcessor(cooker cooker.CookingService, encoder encoder.EncodingService, uploader uploader.UploadingService, progressCh chan processing_common.Watchable, store job_store.JobStore) *RecordProcessor {
	return &RecordProcessor{
		cooker:     cooker,
		encoder:    encoder,
		uploader:   uploader,
		progressCh: progressCh,
		store:      store,
	}
}

// If the system restart for any reason, load all jobs from the store and restart processing them
// They should be able to recover from any state
func (rp *RecordProcessor) Init() error {
	jobs, err := rp.store.GetAll(true)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		go rp.startJob(job)
	}

	return nil
}

func (rp *RecordProcessor) startJob(job *job_store.JobState) {
	err := rp.Process(job)
	if err != nil {
		slog.Error(fmt.Sprintf("[RecordProcessor] :: while processing job %s : %s", job.Id, err.Error()))
		err = rp.store.Delete(job.Id)
		if err != nil {
			slog.Warn(fmt.Sprintf("[RecordProcessor] :: while deleting job %s : %s", job.Id, err.Error()))
		}
		// TODO :: Handle error
	}
}

func (rp *RecordProcessor) Add(audioKeys []string, backgroundAudioKey string) (string, error) {

	jobId := uuid.New().String()
	job := &job_store.JobState{
		Id:                 jobId,
		Step:               processing_common.StepCooking,
		RawAudioKeys:       audioKeys,
		BackgroundAudioKey: backgroundAudioKey,
		UserInput: processing_common.UserInput{Vid: processing_common.VideoOpt{
			Description: "",
			Title:       jobId,
			Visibility:  processing_common.Unlisted,
		}},
	}
	err := rp.store.Upsert(job)
	if err != nil {
		return "", err
	}
	// The channel is buffered for 100 events, it may block if we get (and we will) go over this limit
	// We should probably have a separate goroutine that reads from the channel and sends the events to the client
	go rp.startJob(job)

	return job.Id, nil
}

func (rp *RecordProcessor) Process(job *job_store.JobState) error {
	// Cook
	if job.Step == processing_common.StepCooking {
		slog.Info(fmt.Sprintf("[RecordProcessor] :: Cooking job %s", job.Id))
		err := rp.cook(job)
		if err != nil {
			return err
		}
	}
	slog.Info(fmt.Sprintf("[RecordProcessor] :: Finished cooking job %s", job.Id))

	// Encode
	if job.Step == processing_common.StepEncoding {
		slog.Info(fmt.Sprintf("[RecordProcessor] :: Encoding job %s", job.Id))
		err := rp.encode(job)
		if err != nil {
			return err
		}
	}
	slog.Info(fmt.Sprintf("[RecordProcessor] :: Finished encoding job %s", job.Id))

	// Upload
	if job.Step == processing_common.StepUploading {
		slog.Info(fmt.Sprintf("[RecordProcessor] :: Uploading job %s", job.Id))
		err := rp.upload(job)
		if err != nil {
			return err
		}
	}
	slog.Info(fmt.Sprintf("[RecordProcessor] :: Finished uploading job %s", job.Id))

	if job.Step != processing_common.StepDone {
		return fmt.Errorf("job %s has a unsupported Step '%d'", job.Id, job.Step)
	}

	rp.progressCh <- &processing_common.DoneEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: job.Id,
			State: processing_common.Done,
		},
		Data: processing_common.DoneData{
			Link: job.VideoLink,
		},
	}
	slog.Info(fmt.Sprintf("[RecordProcessor] :: Processing done for job %s", job.Id))

	return nil
}

func (rp *RecordProcessor) UpdateInfos(jobId string, userInput processing_common.UserInput) error {
	job, err := rp.store.Get(jobId)
	if err != nil {
		return err
	}
	job.UserInput = userInput
	err = rp.store.Upsert(job)
	if err != nil {
		return err
	}
	return nil
}

func (rp *RecordProcessor) cook(job *job_store.JobState) error {
	cookedKeys, err := rp.cooker.Cook(job.Id, job.RawAudioKeys)
	if err != nil {
		return err
	}
	job.CookedAudioKeys = cookedKeys
	job.Step = processing_common.StepEncoding
	err = rp.store.Upsert(job)
	if err != nil {
		slog.Warn(fmt.Sprintf("[RecordProcessor] :: while saving job map: %s", err.Error()))
	}
	return nil
}

func (rp *RecordProcessor) encode(job *job_store.JobState) error {
	videoKey, err := rp.encoder.Encode(job.Id, job.CookedAudioKeys, job.BackgroundAudioKey)
	if err != nil {
		return err
	}
	job.VideoKey = videoKey
	job.Step = processing_common.StepUploading
	err = rp.store.Upsert(job)
	if err != nil {
		slog.Warn(fmt.Sprintf("[RecordProcessor] :: while saving job map: %s", err.Error()))
	}
	return nil
}

func (rp *RecordProcessor) upload(job *job_store.JobState) error {
	vid, err := rp.uploader.Upload(job.Id, job.VideoKey, &processing_common.VideoOpt{
		Description: job.UserInput.Vid.Description,
		Title:       job.UserInput.Vid.Title,
		Visibility:  job.UserInput.Vid.Visibility,
	})
	if err != nil {
		return err
	}

	job.VideoLink = vid.WatchPrefix + vid.Id
	job.Step = processing_common.StepDone
	err = rp.store.Upsert(job)
	if err != nil {
		slog.Warn(fmt.Sprintf("[RecordProcessor] :: while saving job map: %s", err.Error()))
	}
	return nil
}

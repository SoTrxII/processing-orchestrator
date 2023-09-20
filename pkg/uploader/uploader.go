package uploader

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"processing-orchestrator/internal/utils"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type Uploader struct {
	invoker         utils.Invoker
	invokeComponent string
	subComponent    string
	// Events received from the Cooking Server
	events chan UploadEvent
	// Channel to send progress to
	progressCh chan processing_common.Watchable
}

func NewUploader(invoker utils.Invoker, invokeComponent, subComponent string, progressCh chan processing_common.Watchable) *Uploader {
	return &Uploader{
		invoker:         invoker,
		invokeComponent: invokeComponent,
		subComponent:    subComponent,
		events:          make(chan UploadEvent, 50),
		progressCh:      progressCh,
	}
}
func (c *Uploader) SubscribeTo(subServer utils.Subscriber) error {
	err := subServer.AddTopicEventHandler(&common.Subscription{
		PubsubName: c.subComponent,
		Topic:      s_Info,
	}, c.onInfo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Uploader) onInfo(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var evt UploadEvent
	err = json.Unmarshal(e.RawData, &evt)
	if err != nil {
		return false, err
	}
	c.events <- evt
	return false, nil
}
func (u *Uploader) handleEvents(ctx context.Context, jobId string) {
Loop:
	for {
		select {
		case evt := <-u.events:
			if evt.JobId != jobId {
				continue
			}
			// The "Done event" from upstream is not propagated
			// for this one, as we will send our own
			/// forged "Done event" when the upload is finished
			if evt.State != processing_common.Done {
				u.progressCh <- &evt
			}

			if evt.State == processing_common.Error || evt.State == processing_common.Done {
				break Loop
			}
		case <-ctx.Done():
		}
	}
}

func (u *Uploader) Upload(jobId, storageKey string, opt *VideoOpt) (*Video, error) {
	ctx, cancel := context.WithCancel(context.Background())
	// So even if an error happen, the goroutine will be cleaned up
	defer cancel()
	go u.handleEvents(ctx, jobId)

	uploadJob := UploadJob{
		VideoOpt:   *opt,
		StorageKey: storageKey,
		JobId:      jobId,
	}
	bytes, err := json.Marshal(uploadJob)
	if err != nil {
		return nil, err
	}
	res, err := u.invoker.InvokeMethodWithContent(context.Background(), u.invokeComponent, "v1/videos", "POST", &utils.DataContent{
		ContentType: "application/json",
		Data:        bytes,
	})
	if err != nil {
		return nil, err
	}
	var video Video
	err = json.Unmarshal(res, &video)
	if err != nil {
		return nil, err
	}
	// Custom done event with video properties
	u.progressCh <- &UploadEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.Done,
		},
		Data: UploadData{
			Link: video.WatchPrefix + video.Id,
		},
	}
	return &video, nil
}

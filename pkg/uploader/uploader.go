package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"processing-orchestrator/internal/demux"
	"processing-orchestrator/internal/utils"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type Uploader struct {
	invoker         utils.Invoker
	invokeComponent string
	subComponent    string
	// Events received from the video store, routed to the job they belong to
	events *demux.Demux[UploadEvent]
	// Channel to send progress to
	progressCh chan processing_common.Watchable
}

func NewUploader(invoker utils.Invoker, invokeComponent, subComponent string, progressCh chan processing_common.Watchable) *Uploader {
	return &Uploader{
		invoker:         invoker,
		invokeComponent: invokeComponent,
		subComponent:    subComponent,
		events:          demux.New[UploadEvent](),
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
	c.events.Dispatch(evt.JobId, evt)
	return false, nil
}
func (u *Uploader) handleEvents(ctx context.Context, events <-chan UploadEvent) {
Loop:
	for {
		select {
		case evt := <-events:
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
			break Loop
		}
	}
}

func (u *Uploader) Upload(jobId, storageKey string, opt *processing_common.VideoOpt) (*Video, error) {
	// Claim the events of this job before asking for the work to be done, so
	// that none of them can be missed
	events, release := u.events.Register(jobId)
	defer release()

	ctx, cancel := context.WithCancel(context.Background())
	// So even if an error happen, the goroutine will be cleaned up
	defer cancel()
	go u.handleEvents(ctx, events)

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
	return &video, nil
}

func (u *Uploader) CreatePlaylist(opt *processing_common.VideoOpt) (*Playlist, error) {
	bytes, err := json.Marshal(opt)
	res, err := u.invoker.InvokeMethodWithContent(context.Background(), u.invokeComponent, "v1/playlists", "POST", &utils.DataContent{
		ContentType: "application/json",
		Data:        bytes,
	})
	if err != nil {
		return nil, err
	}
	var playlist Playlist
	err = json.Unmarshal(res, &playlist)
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (u *Uploader) AddToPlaylist(vidId, playlistId string) error {
	_, err := u.invoker.InvokeMethod(context.Background(), u.invokeComponent, fmt.Sprintf("v1/playlists/%s/videos/%s", playlistId, vidId), "PUT")
	return err
}

func (u *Uploader) SetThumbnail(vidId, thumbKey string) error {
	_, err := u.invoker.InvokeMethod(context.Background(), u.invokeComponent, fmt.Sprintf("v1/videos/%s/thumbnail/%s", vidId, thumbKey), "POST")
	return err
}

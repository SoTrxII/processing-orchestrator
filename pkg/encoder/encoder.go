package encoder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"processing-orchestrator/internal/demux"
	"processing-orchestrator/internal/utils"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type Encoder struct {
	pubClient    utils.Publisher
	pubComponent string
	subComponent string
	// Events received from the encoding service, routed to the job they belong to
	events *demux.Demux[EncodingEvent]
	// Channel to send progress to
	progressCh chan processing_common.Watchable
	opt        *EncoderOpt
}

func NewEncoder(pubClient utils.Publisher, pubComponent, subComponent string, progressCh chan processing_common.Watchable) *Encoder {

	return &Encoder{
		pubClient:    pubClient,
		pubComponent: pubComponent,
		subComponent: subComponent,
		events:       demux.New[EncodingEvent](),
		progressCh:   progressCh,
		opt: &EncoderOpt{
			Extension: ".mp4",
		},
	}
}

func (c *Encoder) SubscribeTo(subServer utils.Subscriber) error {
	err := subServer.AddTopicEventHandler(&common.Subscription{
		PubsubName: c.subComponent,
		Topic:      s_Info,
	}, c.onInfo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Encoder) onInfo(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var evt EncodingEvent
	err = json.Unmarshal(e.RawData, &evt)
	if err != nil {
		return false, err
	}
	c.events.Dispatch(evt.JobId, evt)
	return false, nil
}

func (c *Encoder) Encode(jobId string, audioKeys []string, bgAudioKey string) (string, error) {
	// Claim the events of this job before asking for the work to be done, so
	// that none of them can be missed
	events, release := c.events.Register(jobId)
	defer release()

	err := c.pubClient.PublishEvent(context.Background(), c.pubComponent, p_Start, EncodeJob{
		JobId:              jobId,
		AudiosKeys:         audioKeys,
		BackgroundAudioKey: bgAudioKey,
		Opt:                EncodeJobOpt{deleteAssetsFromObjStore: true},
	})
	if err != nil {
		return "", err
	}

	// TODO : Get this from the encoder itself
	extension := c.opt.Extension
	for {
		select {
		case evt := <-events:
			c.progressCh <- &evt
			switch evt.State {
			case processing_common.Done:
				return jobId + extension, nil
			case processing_common.Error:
				return "", fmt.Errorf(evt.Data.Message)
			}
		}
	}
}

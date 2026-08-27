package cooker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"processing-orchestrator/internal/demux"
	"processing-orchestrator/internal/utils"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type Cooker struct {
	pubClient    utils.Publisher
	pubComponent string
	subComponent string
	// Events received from the Cooking Server, routed to the job they belong to
	events *demux.Demux[CookingEvent]
	// Channel to send progress to
	progressCh chan processing_common.Watchable
	opt        *CookerOpt
}

func NewCooker(pubClient utils.Publisher, pubComponent, subComponent string, progressCh chan processing_common.Watchable) *Cooker {

	return &Cooker{
		pubClient:    pubClient,
		pubComponent: pubComponent,
		subComponent: subComponent,
		events:       demux.New[CookingEvent](),
		progressCh:   progressCh,
		opt: &CookerOpt{
			Extension: ".ogg",
		},
	}
}

func (c *Cooker) SubscribeTo(subServer utils.Subscriber) error {
	err := subServer.AddTopicEventHandler(&common.Subscription{
		PubsubName: c.subComponent,
		Topic:      s_Info,
	}, c.onInfo)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cooker) onInfo(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	var evt CookingEvent
	err = json.Unmarshal(e.RawData, &evt)
	if err != nil {
		return false, err
	}
	c.events.Dispatch(evt.JobId, evt)
	return false, nil
}

func (c *Cooker) Cook(jobId string, recordIds []string) ([]string, error) {
	// Claim the events of this job before asking for the work to be done, so
	// that none of them can be missed
	events, release := c.events.Register(jobId)
	defer release()

	err := c.pubClient.PublishEvent(context.Background(), c.pubComponent, p_Start, CookingJob{
		JobId: jobId,
		Ids:   recordIds,
	})
	if err != nil {
		return nil, err
	}

	doneCh := make(chan error, 1)
	defer close(doneCh)
	var doneEvt CookingEvent
	count := 0
	for {
		select {
		// All records have finished cooking
		case err, _ := <-doneCh:
			if err != nil {
				return nil, err
			}
			ext := c.opt.Extension
			if doneEvt.Data.Extension != "" {
				ext = doneEvt.Data.Extension
			}
			var cookedKeys []string
			for _, id := range recordIds {
				cookedKeys = append(cookedKeys, id+ext)
			}
			return cookedKeys, nil

		// A new event is received
		case evt, _ := <-events:
			c.progressCh <- &evt
			// That can be a done event
			switch evt.State {
			case processing_common.Done:
				if has(recordIds, evt.RecordId) {
					count++
				}
				if count == len(recordIds) {
					doneEvt = evt
					doneCh <- nil
				}
			case processing_common.Error:
				return nil, fmt.Errorf(evt.Data.Message)
			}

		}
	}
}

func has(haystack []string, needle string) bool {
	for _, elem := range haystack {
		if elem == needle {
			return true
		}
	}
	return false
}

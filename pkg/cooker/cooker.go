package cooker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dapr/go-sdk/service/common"
	"processing-orchestrator/internal/utils"
	processing_common "processing-orchestrator/pkg/processing-common"
)

type Cooker struct {
	pubClient    utils.Publisher
	pubComponent string
	subComponent string
	// Events received from the Cooking Server
	events chan CookingEvent
	// Channel to send progress to
	progressCh chan processing_common.Watchable
	opt        *CookerOpt
}

func NewCooker(pubClient utils.Publisher, pubComponent, subComponent string, progressCh chan processing_common.Watchable) *Cooker {

	return &Cooker{
		pubClient:    pubClient,
		pubComponent: pubComponent,
		subComponent: subComponent,
		events:       make(chan CookingEvent, 50),
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
	c.events <- evt
	return false, nil
}

func (c *Cooker) Cook(jobId string, recordIds []string) ([]string, error) {
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
		case evt, _ := <-c.events:
			if evt.JobId != jobId {
				continue
			}
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

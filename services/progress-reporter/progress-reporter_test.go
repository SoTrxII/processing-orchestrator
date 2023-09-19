package progress_reporter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"processing-orchestrator/pkg/cooker"
	processing_common "processing-orchestrator/pkg/processing-common"
	"testing"
)

func TestProgressReporter_ToStatus(t *testing.T) {
	rep := NewProgressReporter()

	// Job with no infos
	status := rep.toStatus(&processing_common.ServiceProgress{
		JobId:       "test",
		Step:        processing_common.StepCooking,
		CurrentItem: "0",
		Error:       nil,
		Progress:    "test",
	})
	assert.Equal(t, status.Id, "test")
	assert.Equal(t, status.CurrentStepIndex, uint32(processing_common.StepCooking))
	// So itemlist is nil
	assert.Nil(t, status.ItemList)
	assert.Equal(t, status.CurrentItemIndex, uint32(0))

	// Adding infos, and current item
	rep.AddInfo("test", []string{"a", "b", "c"})
	status = rep.toStatus(&processing_common.ServiceProgress{
		JobId:       "test",
		Step:        processing_common.StepCooking,
		CurrentItem: "b",
		Error:       nil,
		Progress:    "test",
	})
	assert.Equal(t, status.Id, "test")
	assert.Equal(t, status.CurrentStepIndex, uint32(processing_common.StepCooking))
	assert.Equal(t, status.ItemList, []string{"a", "b", "c"})
	assert.Equal(t, status.CurrentItemIndex, uint32(1))

	// Error
	status = rep.toStatus(&processing_common.ServiceProgress{
		JobId:       "test",
		Step:        processing_common.StepCooking,
		CurrentItem: "b",
		Error:       fmt.Errorf("test"),
		Progress:    "test",
	})
	assert.Equal(t, status.Error, "test")

	// Done
	status = rep.toStatus(&processing_common.ServiceProgress{
		JobId:       "test",
		Step:        processing_common.StepDone,
		CurrentItem: "b",
		Error:       nil,
		Progress:    "test",
	})
	assert.True(t, status.Done)

}

func TestProgressReporter_Start(t *testing.T) {
	rep := NewProgressReporter()
	evtCh := make(chan processing_common.Watchable, 1)
	go rep.Start(evtCh)
	evtCh <- &cooker.CookingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: "test",
			State: processing_common.Done,
		},
		RecordId: "b",
		Data: cooker.CookingData{
			Extension: ".mp3",
		},
	}

}

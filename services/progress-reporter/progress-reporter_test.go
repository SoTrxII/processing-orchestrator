package progress_reporter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"processing-orchestrator/pkg/cooker"
	processing_common "processing-orchestrator/pkg/processing-common"
	pb "processing-orchestrator/proto"
	"sync"
	"testing"
	"time"
)

func TestProgressReporter_ToStatus(t *testing.T) {
	evtCh := make(chan processing_common.Watchable, 1)
	rep := NewProgressReporter(evtCh)

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
	evtCh := make(chan processing_common.Watchable, 1)
	rep := NewProgressReporter(evtCh)
	go rep.Start()
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

func TestProgressReporter_DropsAClientThatStoppedWatching(t *testing.T) {
	evtCh := make(chan processing_common.Watchable, 1)
	rep := NewProgressReporter(evtCh)

	com := &BidirectionalCom{
		Data: make(chan *pb.ProcessingStatus, 1),
		Done: make(chan struct{}),
	}
	rep.Register("test", com)
	// The client is gone. Writing to a channel it closed on its way out used
	// to take the whole process down
	close(com.Done)

	rep.broadcast("test", &pb.ProcessingStatus{Id: "test"})

	assert.Len(t, com.Data, 0)
	rep.mu.Lock()
	defer rep.mu.Unlock()
	assert.NotContains(t, rep.jobWatchers, "test")
}

// Clients register from their own gRPC handler while the reporting loop reads
// and prunes them, and each of them only ever wants its own job
func TestProgressReporter_ConcurrentJobs(t *testing.T) {
	evtCh := make(chan processing_common.Watchable, 100)
	rep := NewProgressReporter(evtCh)
	go rep.Start()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		jobId := fmt.Sprintf("job-%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			com := &BidirectionalCom{
				Data: make(chan *pb.ProcessingStatus, 10),
				Done: make(chan struct{}),
			}
			defer close(com.Done)
			rep.Register(jobId, com)
			rep.AddInfo(jobId, []string{"a"})

			evtCh <- &cooker.CookingEvent{
				ServiceEvent: processing_common.ServiceEvent{
					JobId: jobId,
					State: processing_common.InProgress,
				},
				RecordId: "a",
			}

			select {
			case <-time.After(5 * time.Second):
				t.Errorf("job %s never got its progress", jobId)
			case status := <-com.Data:
				assert.Equal(t, jobId, status.Id)
			}
		}()
	}
	wg.Wait()
}

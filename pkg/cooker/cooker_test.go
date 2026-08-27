package cooker

import (
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	processing_common "processing-orchestrator/pkg/processing-common"
	test_utils "processing-orchestrator/test-utils"
	"testing"
	"time"
)

func setup() *Cooker {
	cooker, _ := setupWithProgress()
	return cooker
}

// setupWithProgress also hands back the channel the cooker reports progress on
func setupWithProgress() (*Cooker, chan processing_common.Watchable) {
	progressCh := make(chan processing_common.Watchable, 100)
	pub := test_utils.MockPublisher{}
	pub.On("PublishEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return NewCooker(&pub, "test", "test", progressCh), progressCh
}

func TestCooker_Start_OneRecord_CustomExt(t *testing.T) {
	cooker := setup()
	done := make(chan bool)
	defer close(done)

	// Cook cooking
	go func() {
		audioKey, err := cooker.Cook("1", []string{"1"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"1.test"}, audioKey)
		done <- true
	}()

	waitUntilWaiting(t, cooker, "1")

	// Record has finished cooking
	data := getDoneEvt("1", "1", ".test")
	_, err := cooker.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// Now the cooking process should be finished
	select {
	case <-time.After(5 * time.Second):
		t.Fail()
		break
	case <-done:
	}

}

func TestCooker_Start_MultipleRecords_CustomExt(t *testing.T) {
	cooker := setup()
	done := make(chan bool)
	defer close(done)

	// Cook cooking
	go func() {
		audioKey, err := cooker.Cook("1", []string{"1", "2"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"1.ogg", "2.ogg"}, audioKey)
		done <- true
	}()

	waitUntilWaiting(t, cooker, "1")

	// First record has finished cooking but the second one is still cooking
	data := getDoneEvt("1", "1", "")
	_, err := cooker.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// So the cooking process should not be finished yet
	select {
	case <-done:
		t.Errorf("Cooking stopped after only one record was processed")
		t.Fail()
	default:
	}

	// Second record has finished cooking
	data = getDoneEvt("1", "2", "")
	_, err = cooker.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// Now the cooking process should be finished
	select {
	case <-time.After(5 * time.Second):
		t.Fail()
		break
	case <-done:
	}

}

func TestCooker_Start_ErrorDuringCooking(t *testing.T) {
	cooker := setup()
	done := make(chan bool)
	defer close(done)

	// Cook cooking
	go func() {
		_, err := cooker.Cook("1", []string{"1", "2"})
		assert.Error(t, err)
		done <- true
	}()

	waitUntilWaiting(t, cooker, "1")
	// Generate error
	data := getErrorEvt("1", "1", "Test Error")
	_, err := cooker.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// This will trigger the error and enc the cooking process
	<-done
}

func getErrorEvt(jobId, recId, msg string) []byte {
	evt := CookingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.Error,
		},
		RecordId: recId,
		Data: CookingData{
			TotalBytes: 0,
			Message:    msg,
			Extension:  "",
		},
	}
	rawEvt, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return rawEvt
}
func getDoneEvt(jobId, recId, ext string) []byte {
	evt := CookingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.Done,
		},
		RecordId: recId,
		Data: CookingData{
			TotalBytes: 0,
			Message:    "",
			Extension:  ext,
		},
	}
	rawEvt, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return rawEvt
}

// waitUntilWaiting blocks until the job started waiting for its events, which
// only happens once the Cook goroutine has been scheduled. Events published
// before that point belong to nobody and are dropped
func waitUntilWaiting(t *testing.T, cooker *Cooker, jobId string) {
	assert.Eventually(t, func() bool {
		return cooker.events.IsWaiting(jobId)
	}, 5*time.Second, 10*time.Millisecond)
}

// Two jobs cooking at the same time each have to get their own events. They
// all come in on the same subscription, and one handed to the wrong job would
// be consumed for good, leaving the right one waiting forever
func TestCooker_Start_ConcurrentJobs(t *testing.T) {
	cooker, progressCh := setupWithProgress()
	first := make(chan []string, 1)
	second := make(chan []string, 1)

	go func() {
		cookedKeys, err := cooker.Cook("job-1", []string{"1"})
		assert.NoError(t, err)
		first <- cookedKeys
	}()
	go func() {
		cookedKeys, err := cooker.Cook("job-2", []string{"2"})
		assert.NoError(t, err)
		second <- cookedKeys
	}()

	waitUntilWaiting(t, cooker, "job-1")
	waitUntilWaiting(t, cooker, "job-2")

	// Every event of the second job has to reach it : one swallowed by the
	// first job on the way would be gone for good
	for i := 0; i < 10; i++ {
		_, err := cooker.onInfo(nil, &common.TopicEvent{RawData: getProgressEvt("job-2", "2", int64(i))})
		assert.NoError(t, err)
	}

	// Finishing the second one first, so a job picking up whatever arrives
	// next would be caught out
	_, err := cooker.onInfo(nil, &common.TopicEvent{RawData: getDoneEvt("job-2", "2", ".ogg")})
	assert.NoError(t, err)
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("the second job never got its own done event")
	case cookedKeys := <-second:
		assert.Equal(t, []string{"2.ogg"}, cookedKeys)
	}

	// The 10 progress events plus the done one, all of them the second job
	assert.Len(t, progressCh, 11)
	for len(progressCh) > 0 {
		assert.Equal(t, "job-2", (<-progressCh).ToProgress().JobId)
	}

	// The first one is still cooking, untouched by the other job
	assert.Len(t, first, 0)
	_, err = cooker.onInfo(nil, &common.TopicEvent{RawData: getDoneEvt("job-1", "1", ".ogg")})
	assert.NoError(t, err)
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("the first job never got its own done event")
	case cookedKeys := <-first:
		assert.Equal(t, []string{"1.ogg"}, cookedKeys)
	}
}

func getProgressEvt(jobId, recId string, totalBytes int64) []byte {
	evt := CookingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.InProgress,
		},
		RecordId: recId,
		Data: CookingData{
			TotalBytes: totalBytes,
		},
	}
	rawEvt, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return rawEvt
}

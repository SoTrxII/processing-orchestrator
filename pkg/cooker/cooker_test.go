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
	progressCh := make(chan processing_common.Watchable, 100)
	pub := test_utils.MockPublisher{}
	sub := test_utils.MockSubscriber{}
	sub.On("AddTopicEventHandler", mock.Anything, mock.Anything).Return(nil)
	pub.On("PublishEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return NewCooker(&pub, &sub, "test", progressCh)
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

package encoder

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

func setup(t *testing.T) *Encoder {
	progressCh := make(chan processing_common.Watchable, 100)
	pub := test_utils.MockPublisher{}
	sub := test_utils.MockSubscriber{}
	sub.On("AddTopicEventHandler", mock.Anything, mock.Anything).Return(nil)
	pub.On("PublishEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return NewEncoder(&pub, "test", "test", progressCh)
}

func TestEncoder_Start(t *testing.T) {
	encoder := setup(t)
	done := make(chan bool)
	defer close(done)

	// Cook encoding
	go func() {
		audioKey, err := encoder.Encode("1", []string{"1"}, "")
		assert.NoError(t, err)
		assert.Equal(t, "1.mp4", audioKey)
		done <- true
	}()

	// Record has finished encoding
	data := getDoneEvt("1")
	_, err := encoder.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// Now the encoding process should be finished
	select {
	case <-time.After(5 * time.Second):
		t.Fail()
		break
	case <-done:

	}

}

func TestEncoder_ErrorDuringEncoding(t *testing.T) {
	encoder := setup(t)
	done := make(chan bool)
	defer close(done)

	// Cook encoding
	go func() {
		videoKey, err := encoder.Encode("1", []string{"1"}, "")
		assert.Error(t, err)
		assert.Empty(t, videoKey)
		done <- true
	}()

	// Record has encountered an error during encoding
	data := getErrorEvt("1")
	_, err := encoder.onInfo(nil, &common.TopicEvent{RawData: data})
	assert.NoError(t, err)

	// Now the encoding process should be finished
	select {
	case <-time.After(5 * time.Second):
		t.Fail()
		break
	case <-done:

	}

}

func getDoneEvt(jobId string) []byte {
	evt := EncodingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.Done,
		},
		Data: EncodingData{
			Frames:  4545421,
			Fps:     60,
			Quality: 20,
			Size:    8898,
			Time:    time.Now(),
			Bitrate: "192kb",
			Speed:   2,
		},
	}
	rawEvt, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return rawEvt
}

func getErrorEvt(jobId string) []byte {
	evt := EncodingEvent{
		ServiceEvent: processing_common.ServiceEvent{
			JobId: jobId,
			State: processing_common.Error,
		},
		Data: EncodingData{
			Message: "Error",
		},
	}
	rawEvt, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	return rawEvt
}

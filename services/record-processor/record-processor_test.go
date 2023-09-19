package record_processor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	test_utils "processing-orchestrator/test-utils"
	"testing"
)

func setup(t *testing.T) *RecordProcessor {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	mockStore := &test_utils.MockJobStore{}
	return NewRecordProcessor(mockCooker, mockEncoder, mockUploader, mockStore)
}

func TestRecordProcessor_Process_ErrorDuringCooking(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, mockStore)

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	assert.Error(t, err)
}

func TestRecordProcessor_Process_ErrorDuringEncoding(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, mockStore)

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	assert.Error(t, err)
}

func TestRecordProcessor_Process_ErrorDuringUploading(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, mockStore)

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	assert.Error(t, err)
}

func TestRecordProcessor_Process(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, mockStore)

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	assert.NoError(t, err)
}

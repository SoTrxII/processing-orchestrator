package record_processor

import (
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

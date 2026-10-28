package record_processor

import (
	"context"
	"fmt"
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	test_utils "processing-orchestrator/test-utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

// A job read back from the store, after a restart, still asks for its work in
// the trace it was started in
func TestAStoredJobWorksInTheTraceItWasStartedIn(t *testing.T) {
	tp := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	mockCooker := test_utils.NewMockCookingService(t)
	var sent []string
	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, _ string, _ []string) {
			md, _ := metadata.FromOutgoingContext(ctx)
			sent = md.Get("traceparent")
		}).
		Return(nil, fmt.Errorf("stop here"))

	rp := NewRecordProcessor(mockCooker, nil, nil, make(chan processing_common.Watchable, 1), &test_utils.MockJobStore{}, Addons{})
	_ = rp.Process(&job_store.JobState{Id: "job-1", Step: processing_common.StepCooking, TraceParent: tp})

	assert.Equal(t, []string{tp}, sent)
}

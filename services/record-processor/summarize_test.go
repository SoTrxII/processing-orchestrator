package record_processor

import (
	"fmt"
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	"processing-orchestrator/pkg/summarizer"
	test_utils "processing-orchestrator/test-utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// A cooked job, stopped right after cooking so the summarizer is the only
// thing under test.
func cookedJob(summary *processing_common.SummaryOpt) *job_store.JobState {
	return &job_store.JobState{
		Id:           "job-1",
		Step:         processing_common.StepCooking,
		RawAudioKeys: []string{"raw"},
		UserInput:    processing_common.UserInput{Summary: summary},
	}
}

func cookOnly(t *testing.T, addons Addons, cooked []string, job *job_store.JobState) error {
	t.Helper()
	mockCooker := &test_utils.MockCookingService{}
	mockStore := &test_utils.MockJobStore{}
	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return(cooked, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)

	rp := NewRecordProcessor(mockCooker, nil, nil, make(chan processing_common.Watchable, 1), mockStore, addons)
	return rp.cook(job)
}

// The summary is submitted as soon as the audio exists, not at the end of the
// job : transcription is the longest thing in the whole pipeline.
func TestSummaryIsSubmittedOnceCookingIsDone(t *testing.T) {
	mockSummarizer := test_utils.NewMockSummarizingService(t)
	var got *summarizer.SummaryJob
	mockSummarizer.EXPECT().Submit(mock.Anything).
		Run(func(req *summarizer.SummaryJob) { got = req }).
		Return(nil)

	job := cookedJob(&processing_common.SummaryOpt{CampaignId: 28, EpisodeId: 11, IsOneShot: false})
	err := cookOnly(t, Addons{Summarizer: mockSummarizer}, []string{"cooked.ogg"}, job)

	assert.NoError(t, err)
	assert.Equal(t, &summarizer.SummaryJob{
		JobId:      "job-1",
		CampaignId: 28,
		EpisodeId:  11,
		IsOneShot:  false,
		AudioKeys:  []string{"cooked.ogg"},
	}, got)
}

// A recording that dropped and resumed is cooked into several files. All of
// them are the session, so all of them have to be handed over.
func TestEveryCookedChunkIsSubmitted(t *testing.T) {
	mockSummarizer := test_utils.NewMockSummarizingService(t)
	var got *summarizer.SummaryJob
	mockSummarizer.EXPECT().Submit(mock.Anything).
		Run(func(req *summarizer.SummaryJob) { got = req }).
		Return(nil)

	job := cookedJob(&processing_common.SummaryOpt{CampaignId: 28, EpisodeId: 12})
	err := cookOnly(t, Addons{Summarizer: mockSummarizer}, []string{"a.ogg", "b.ogg"}, job)

	assert.NoError(t, err)
	assert.Equal(t, []string{"a.ogg", "b.ogg"}, got.AudioKeys)
}

func TestOneShotIsFlagged(t *testing.T) {
	mockSummarizer := test_utils.NewMockSummarizingService(t)
	var got *summarizer.SummaryJob
	mockSummarizer.EXPECT().Submit(mock.Anything).
		Run(func(req *summarizer.SummaryJob) { got = req }).
		Return(nil)

	job := cookedJob(&processing_common.SummaryOpt{CampaignId: 3, EpisodeId: 1, IsOneShot: true})
	err := cookOnly(t, Addons{Summarizer: mockSummarizer}, []string{"cooked.ogg"}, job)

	assert.NoError(t, err)
	assert.True(t, got.IsOneShot)
}

// Nobody asked for a summary : the recording is processed as it always was.
func TestNoCampaignMeansNoSummary(t *testing.T) {
	mockSummarizer := test_utils.NewMockSummarizingService(t)

	job := cookedJob(nil)
	err := cookOnly(t, Addons{Summarizer: mockSummarizer}, []string{"cooked.ogg"}, job)

	assert.NoError(t, err)
	mockSummarizer.AssertNotCalled(t, "Submit", mock.Anything)
}

// No summarizer wired at all has to stay a working configuration.
func TestNoSummarizerIsFine(t *testing.T) {
	job := cookedJob(&processing_common.SummaryOpt{CampaignId: 28, EpisodeId: 11})
	err := cookOnly(t, Addons{}, []string{"cooked.ogg"}, job)
	assert.NoError(t, err)
}

// The video is what the user is waiting for. A summarizer that is down must
// not take the whole job with it.
func TestAFailingSummarizerDoesNotFailTheJob(t *testing.T) {
	mockSummarizer := test_utils.NewMockSummarizingService(t)
	mockSummarizer.EXPECT().Submit(mock.Anything).Return(fmt.Errorf("connection refused"))

	job := cookedJob(&processing_common.SummaryOpt{CampaignId: 28, EpisodeId: 11})
	err := cookOnly(t, Addons{Summarizer: mockSummarizer}, []string{"cooked.ogg"}, job)

	assert.NoError(t, err)
	assert.Equal(t, processing_common.StepEncoding, job.Step, "cooking still hands over to encoding")
}

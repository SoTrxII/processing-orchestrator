package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"processing-orchestrator/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// recordingInvoker stands in for the Dapr sidecar and keeps what it was asked
// to send.
type recordingInvoker struct {
	appId   string
	method  string
	verb    string
	content *utils.DataContent
	err     error
}

func (r *recordingInvoker) InvokeMethodWithContent(_ context.Context, appID, method, verb string, content *utils.DataContent) ([]byte, error) {
	r.appId, r.method, r.verb, r.content = appID, method, verb, content
	return nil, r.err
}

func (r *recordingInvoker) InvokeMethod(_ context.Context, appID, method, verb string) ([]byte, error) {
	r.appId, r.method, r.verb = appID, method, verb
	return nil, r.err
}

// The field names here are a wire contract with summary-orchestrator's
// /workflows/episode endpoint : renaming one of them silently stops the
// summary from being filed under the right campaign.
func TestSubmitSpeaksTheEndpointsLanguage(t *testing.T) {
	invoker := &recordingInvoker{}
	s := NewSummarizer(invoker, "summary-orchestrator")

	err := s.Submit(&SummaryJob{
		JobId:      "job-1",
		CampaignId: 28,
		EpisodeId:  11,
		IsOneShot:  true,
		AudioKeys:  []string{"a.ogg", "b.ogg"},
	})
	assert.NoError(t, err)

	assert.Equal(t, "summary-orchestrator", invoker.appId)
	assert.Equal(t, "workflows/episode", invoker.method)
	assert.Equal(t, "POST", invoker.verb)
	assert.Equal(t, "application/json", invoker.content.ContentType)

	var sent map[string]any
	assert.NoError(t, json.Unmarshal(invoker.content.Data, &sent))
	assert.Equal(t, map[string]any{
		"jobId":      "job-1",
		"campaignId": float64(28),
		"episodeId":  float64(11),
		"isOneShot":  true,
		"audioKeys":  []any{"a.ogg", "b.ogg"},
	}, sent)
}

func TestSubmitReportsAnUnreachableSummarizer(t *testing.T) {
	invoker := &recordingInvoker{err: fmt.Errorf("connection refused")}
	s := NewSummarizer(invoker, "summary-orchestrator")

	err := s.Submit(&SummaryJob{JobId: "job-1"})
	assert.ErrorContains(t, err, "job-1")
	assert.ErrorContains(t, err, "connection refused")
}

// Empty must not appear on the wire at all: a recording with no game master
// recorded has to produce the same request as one sent before the field
// existed, so an older summary-orchestrator keeps working during a rollout.
func TestSubmitOmitsGameMastersWhenThereAreNone(t *testing.T) {
	invoker := &recordingInvoker{}
	s := NewSummarizer(invoker, "summary-orchestrator")

	assert.NoError(t, s.Submit(&SummaryJob{JobId: "job-1", CampaignId: 28, EpisodeId: 11}))

	var sent map[string]any
	assert.NoError(t, json.Unmarshal(invoker.content.Data, &sent))
	_, present := sent["gameMasterIds"]
	assert.False(t, present, "an empty game master list must not be serialized")
}

// And must travel when there is one: this is the whole reason the summarizer
// no longer reads Velvet's database itself.
func TestSubmitCarriesGameMasters(t *testing.T) {
	invoker := &recordingInvoker{}
	s := NewSummarizer(invoker, "summary-orchestrator")

	assert.NoError(t, s.Submit(&SummaryJob{
		JobId:         "job-1",
		CampaignId:    28,
		EpisodeId:     11,
		GameMasterIds: []string{"188626510901542912"},
	}))

	var sent map[string]any
	assert.NoError(t, json.Unmarshal(invoker.content.Data, &sent))
	assert.Equal(t, []any{"188626510901542912"}, sent["gameMasterIds"])
}

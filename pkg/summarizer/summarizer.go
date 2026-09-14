package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"processing-orchestrator/internal/utils"
	"time"
)

// startMethod is the route on summary-orchestrator that schedules a
// summarization.
const startMethod = "workflows/episode"

// submitTimeout only has to cover scheduling the work. The endpoint answers
// as soon as the workflow is queued, so anything longer than this means the
// service is not answering rather than thinking.
const submitTimeout = 30 * time.Second

// Summarizer reaches summary-orchestrator through the Dapr sidecar.
type Summarizer struct {
	client utils.Invoker
	appId  string
}

func NewSummarizer(client utils.Invoker, appId string) *Summarizer {
	return &Summarizer{client: client, appId: appId}
}

func (s *Summarizer) Submit(req *SummaryJob) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), submitTimeout)
	defer cancel()

	_, err = s.client.InvokeMethodWithContent(ctx, s.appId, startMethod, "POST", &utils.DataContent{
		ContentType: "application/json",
		Data:        payload,
	})
	if err != nil {
		return fmt.Errorf("submitting job %s to %s: %w", req.JobId, s.appId, err)
	}
	return nil
}

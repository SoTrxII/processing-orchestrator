package summarizer

// SummarizingService hands a finished recording over to be summarized.
//
// It is deliberately fire-and-forget : summarizing a session is hours of
// transcription and model time, far longer than the job that triggers it.
// Submit returns as soon as the work has been accepted, and nothing here ever
// waits for a summary to exist.
type SummarizingService interface {
	// Submit asks for a recording to be summarized. Returning nil means the
	// request was accepted, not that the summary is ready.
	Submit(req *SummaryJob) error
}

// SummaryJob is what the summarizer is given. The field names are the wire
// contract with summary-orchestrator's /workflows/episode endpoint.
type SummaryJob struct {
	// Processing job this recording belongs to. The summarizer uses it as its
	// own workflow id, which is what makes submitting twice a no-op instead
	// of a second transcription.
	JobId string `json:"jobId"`
	// Campaign this episode belongs to
	CampaignId int `json:"campaignId"`
	// Episode number within the campaign
	EpisodeId int `json:"episodeId"`
	// A one-shot has no campaign narrative to rebuild
	IsOneShot bool `json:"isOneShot"`
	// Discord ids of whoever runs this campaign. Lets the summarizer tell a
	// player character from an NPC without querying Velvet's database.
	//
	// omitempty so a caller that supplies none sends exactly the payload it
	// sent before this field existed, which is what keeps a rolling deploy
	// against an older summary-orchestrator byte-identical rather than merely
	// tolerable.
	GameMasterIds []string `json:"gameMasterIds,omitempty"`
	// Keys of the cooked audio, in recording order. A session that dropped
	// and resumed has several.
	AudioKeys []string `json:"audioKeys"`
}

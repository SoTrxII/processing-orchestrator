package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"processing-orchestrator/internal/utils"
	job_store "processing-orchestrator/pkg/job-store"
	"strings"
	"time"

	"github.com/dapr/go-sdk/service/common"
)

// GistTopic is where summary-orchestrator announces a finished summary.
const GistTopic = "gists"

// chaptersTimeout covers one model call on the storyteller's side: about a
// minute measured on seven sessions, up to 50k characters of prompt for a
// five-hour one.
const chaptersTimeout = 10 * time.Minute

// The chapters are kept between these two lines, so that writing them again
// replaces the block instead of stacking a second one. Neither line may look
// like a timestamp, or YouTube would read it as a chapter.
const (
	chaptersStart = "Chapitres :"
	chaptersEnd   = "(chapitres générés automatiquement)"
)

// JobLookup finds a processing job by id.
type JobLookup interface {
	Get(id string) (*job_store.JobState, error)
}

// ChapterPublisher puts YouTube chapters on a session's video once its summary
// exists. Summarizing stays optional: this only runs when the summarizer addon
// does, and a failure here leaves a video without chapters, nothing worse.
type ChapterPublisher struct {
	invoker      utils.Invoker
	summaryAppId string
	videoAppId   string
	pubsub       string
	jobs         JobLookup
}

func NewChapterPublisher(invoker utils.Invoker, summaryAppId, videoAppId, pubsub string, jobs JobLookup) *ChapterPublisher {
	return &ChapterPublisher{invoker, summaryAppId, videoAppId, pubsub, jobs}
}

func (p *ChapterPublisher) SubscribeTo(sub utils.Subscriber) error {
	return sub.AddTopicEventHandler(&common.Subscription{PubsubName: p.pubsub, Topic: GistTopic}, p.onGist)
}

// gist is the part of summary-orchestrator's Gist event used here.
type gist struct {
	JobId      string `json:"jobId"`
	CampaignId int    `json:"campaignId"`
	EpisodeId  int    `json:"episodeId"`
}

// onGist acknowledges at once and does the work aside: it takes a model call,
// and the Redis pubsub redelivers a message left unacknowledged past its
// processing timeout (60 s by default), which would ask for the chapters again.
func (p *ChapterPublisher) onGist(_ context.Context, e *common.TopicEvent) (retry bool, err error) {
	var g gist
	if err := json.Unmarshal(e.RawData, &g); err != nil {
		return false, err
	}
	// ponytail: work in flight is lost if the pod stops during it; re-post the
	// gist or call the chapters route by hand if that ever matters.
	go func() {
		if err := p.Publish(context.Background(), g); err != nil {
			slog.Warn(fmt.Sprintf("[Chapters] :: campaign %d, episode %d (job %s): %s", g.CampaignId, g.EpisodeId, g.JobId, err))
		}
	}()
	return false, nil
}

// chapter is one entry of summary-orchestrator's chapters route.
type chapter struct {
	Start int    `json:"start"`
	Title string `json:"title"`
}

// Publish asks for the episode's chapters and writes them into the
// description of the video its job uploaded.
func (p *ChapterPublisher) Publish(ctx context.Context, g gist) error {
	// A gist with no job is a summary run by hand: there is no video of ours
	// to put the chapters on.
	if g.JobId == "" {
		return nil
	}
	job, err := p.jobs.Get(g.JobId)
	if err != nil {
		return err
	}
	if job == nil || job.VideoLink == "" {
		slog.Info(fmt.Sprintf("[Chapters] :: no uploaded video for job %s, no chapters", g.JobId))
		return nil
	}
	videoId := videoIdOf(job.VideoLink)

	chapters, err := p.chapters(ctx, g)
	if err != nil {
		return err
	}
	if len(chapters) < 3 {
		slog.Info(fmt.Sprintf("[Chapters] :: no chapters for campaign %d, episode %d", g.CampaignId, g.EpisodeId))
		return nil
	}

	raw, err := p.invoker.InvokeMethod(ctx, p.videoAppId, "v1/videos/"+videoId, "GET")
	if err != nil {
		return fmt.Errorf("reading video %s: %w", videoId, err)
	}
	// Sent back whole: video-store's PUT sets title, description and
	// visibility from the body, and refuses a changed id or creation date.
	var video map[string]any
	if err := json.Unmarshal(raw, &video); err != nil {
		return fmt.Errorf("reading video %s: %w", videoId, err)
	}
	current, _ := video["description"].(string)
	updated := withChapters(current, chapterBlock(chapters))
	if updated == current {
		return nil
	}
	video["description"] = updated
	body, _ := json.Marshal(video)
	if _, err := p.invoker.InvokeMethodWithContent(ctx, p.videoAppId, "v1/videos/"+videoId, "PUT",
		&utils.DataContent{ContentType: "application/json", Data: body}); err != nil {
		return fmt.Errorf("updating video %s: %w", videoId, err)
	}
	slog.Info(fmt.Sprintf("[Chapters] :: %d chapters on video %s (campaign %d, episode %d)", len(chapters), videoId, g.CampaignId, g.EpisodeId))
	return nil
}

func (p *ChapterPublisher) chapters(ctx context.Context, g gist) ([]chapter, error) {
	ctx, cancel := context.WithTimeout(ctx, chaptersTimeout)
	defer cancel()
	raw, err := p.invoker.InvokeMethod(ctx, p.summaryAppId,
		fmt.Sprintf("campaigns/%d/episodes/%d/chapters", g.CampaignId, g.EpisodeId), "POST")
	if err != nil {
		return nil, fmt.Errorf("asking for chapters: %w", err)
	}
	var out struct {
		Chapters []chapter `json:"chapters"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("reading chapters: %w", err)
	}
	return out.Chapters, nil
}

// videoIdOf takes the id off a watch link (job.VideoLink is the store's watch
// prefix plus the id).
func videoIdOf(link string) string {
	if i := strings.LastIndex(link, "v="); i >= 0 {
		return link[i+2:]
	}
	return link[strings.LastIndex(link, "/")+1:]
}

// clock writes a chapter start the way YouTube's help shows it: 00:00, and
// h:mm:ss past the hour.
func clock(s int) string {
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s%3600/60, s%60)
	}
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func chapterBlock(chapters []chapter) string {
	lines := []string{chaptersStart}
	for _, c := range chapters {
		lines = append(lines, clock(c.Start)+" "+c.Title)
	}
	return strings.Join(append(lines, chaptersEnd), "\n")
}

// withChapters replaces the chapters block of a description, or appends one.
func withChapters(description, block string) string {
	start := strings.Index(description, chaptersStart)
	end := strings.Index(description, chaptersEnd)
	if start >= 0 && end > start {
		return description[:start] + block + description[end+len(chaptersEnd):]
	}
	if strings.TrimSpace(description) == "" {
		return block
	}
	return strings.TrimRight(description, "\n") + "\n\n" + block
}

package summarizer

import (
	"context"
	"encoding/json"
	"processing-orchestrator/internal/utils"
	job_store "processing-orchestrator/pkg/job-store"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeSidecar answers like summary-orchestrator and video-store would, and
// keeps the description it was asked to write.
type fakeSidecar struct {
	chapters string
	video    map[string]any
	calls    []string
	written  map[string]any
}

func (f *fakeSidecar) InvokeMethod(_ context.Context, appID, method, verb string) ([]byte, error) {
	f.calls = append(f.calls, verb+" "+appID+"/"+method)
	if appID == "summary-orchestrator" {
		return []byte(f.chapters), nil
	}
	return json.Marshal(f.video)
}

func (f *fakeSidecar) InvokeMethodWithContent(_ context.Context, appID, method, verb string, content *utils.DataContent) ([]byte, error) {
	f.calls = append(f.calls, verb+" "+appID+"/"+method)
	_ = json.Unmarshal(content.Data, &f.written)
	return nil, nil
}

type jobs map[string]*job_store.JobState

func (j jobs) Get(id string) (*job_store.JobState, error) { return j[id], nil }

func publisher(f *fakeSidecar, j jobs) *ChapterPublisher {
	return NewChapterPublisher(f, "summary-orchestrator", "video-store", "pubsub", j)
}

const threeChapters = `{"chapters":[{"start":0,"title":"Avant la partie"},{"start":296,"title":"Cachés sous la caverne"},{"start":7658,"title":"Duel pour le fruit"}]}`

func TestChaptersGoInTheVideosDescription(t *testing.T) {
	f := &fakeSidecar{chapters: threeChapters, video: map[string]any{
		"id": "znBAEoKRsYc", "title": "Tigrou - Episode 9", "description": "Bonne séance !",
		"createdAt": "2026-11-02T23:30:00Z", "visibility": "unlisted"}}
	p := publisher(f, jobs{"job-9": {Id: "job-9", VideoLink: "https://www.youtube.com/watch?v=znBAEoKRsYc"}})

	var g gist
	assert.NoError(t, json.Unmarshal([]byte(`{"jobId":"job-9","campaignId":22,"episodeId":9,"summary":"…"}`), &g))
	assert.NoError(t, p.Publish(context.Background(), g))

	assert.Equal(t, []string{
		"POST summary-orchestrator/campaigns/22/episodes/9/chapters",
		"GET video-store/v1/videos/znBAEoKRsYc",
		"PUT video-store/v1/videos/znBAEoKRsYc",
	}, f.calls)
	assert.Equal(t, "Bonne séance !\n\nChapitres :\n00:00 Avant la partie\n04:56 Cachés sous la caverne\n2:07:38 Duel pour le fruit\n(chapitres générés automatiquement)",
		f.written["description"])
	// Everything else goes back as it came: video-store refuses a changed id or date.
	assert.Equal(t, "2026-11-02T23:30:00Z", f.written["createdAt"])
	assert.Equal(t, "unlisted", f.written["visibility"])
}

// Written again (a summary re-run), the block is replaced, not stacked.
func TestChaptersReplaceTheirOwnBlock(t *testing.T) {
	old := "Intro\n\nChapitres :\n00:00 Ancien\n(chapitres générés automatiquement)\nFin"
	got := withChapters(old, chapterBlock([]chapter{{0, "A"}, {60, "B"}, {120, "C"}}))
	assert.Equal(t, "Intro\n\nChapitres :\n00:00 A\n01:00 B\n02:00 C\n(chapitres générés automatiquement)\nFin", got)
	assert.Equal(t, 1, strings.Count(got, chaptersStart))
	assert.Equal(t, chapterBlock([]chapter{{0, "A"}}), withChapters("  ", chapterBlock([]chapter{{0, "A"}})))
}

// Nothing to do, and nothing called on YouTube's quota: a summary run by
// hand (no job), a job whose upload never finished, or too few chapters.
func TestNoChaptersWithoutAVideoOrEnoughChapters(t *testing.T) {
	for name, tc := range map[string]struct {
		gist     gist
		jobs     jobs
		chapters string
	}{
		"no job id":   {gist{CampaignId: 1, EpisodeId: 1}, jobs{}, threeChapters},
		"unknown job": {gist{JobId: "x", CampaignId: 1, EpisodeId: 1}, jobs{}, threeChapters},
		"no video":    {gist{JobId: "j", CampaignId: 1, EpisodeId: 1}, jobs{"j": {Id: "j"}}, threeChapters},
		"too few":     {gist{JobId: "j", CampaignId: 1, EpisodeId: 1}, jobs{"j": {Id: "j", VideoLink: "https://www.youtube.com/watch?v=abc"}}, `{"chapters":[{"start":0,"title":"a"}]}`},
	} {
		f := &fakeSidecar{chapters: tc.chapters}
		assert.NoError(t, publisher(f, tc.jobs).Publish(context.Background(), tc.gist), name)
		for _, c := range f.calls {
			assert.NotContains(t, c, "video-store", name)
		}
	}
}

func TestVideoIdOf(t *testing.T) {
	assert.Equal(t, "znBAEoKRsYc", videoIdOf("https://www.youtube.com/watch?v=znBAEoKRsYc"))
	assert.Equal(t, "pAuH5rMegqI", videoIdOf("https://youtu.be/pAuH5rMegqI"))
}

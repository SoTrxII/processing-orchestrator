package main

import (
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	pb "processing-orchestrator/proto"
	progress_reporter "processing-orchestrator/services/progress-reporter"
	"testing"
)

func TestSettledStatus_AFinishedJobIsDoneWithItsLinks(t *testing.T) {
	rep := progress_reporter.NewProgressReporter(nil)
	job := &job_store.JobState{
		Id:                  "job",
		Step:                processing_common.StepDone,
		VideoLink:           "https://youtu.be/v",
		CreatedPlaylistLink: "https://youtube.com/playlist?list=p",
	}

	status := settledStatus("job", job, rep)
	if status == nil || !status.Done {
		t.Fatalf("a finished job must be reported done, got %+v", status)
	}
	if status.Link != job.VideoLink || status.CreatedPlaylistLink != job.CreatedPlaylistLink {
		t.Errorf("links lost : %+v", status)
	}
	// Velvet formats its progress message from the steps
	if len(status.StepsList) == 0 {
		t.Error("no steps to display")
	}
}

func TestSettledStatus_AnUnknownJobFailed(t *testing.T) {
	status := settledStatus("gone", nil, progress_reporter.NewProgressReporter(nil))
	if status == nil || status.Error == "" || status.Done {
		t.Fatalf("an unknown job must be reported failed, got %+v", status)
	}
}

func TestSettledStatus_ARunningJobIsLeftToItsEvents(t *testing.T) {
	job := &job_store.JobState{Id: "job", Step: processing_common.StepEncoding}
	if status := settledStatus("job", job, progress_reporter.NewProgressReporter(nil)); status != nil {
		t.Fatalf("a running job has nothing settled yet, got %+v", status)
	}
}

func TestUserInput_ANilThumbnailIsAnEmptyOne(t *testing.T) {
	input, err := userInput(&pb.UpdateRequest{Id: "job", VidTitle: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if input.Vid.Title != "t" || input.Vid.Thumbnail != (processing_common.ThumbnailOpt{}) {
		t.Errorf("got %+v", input.Vid)
	}
}

package record_processor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	job_store "processing-orchestrator/pkg/job-store"
	processing_common "processing-orchestrator/pkg/processing-common"
	"processing-orchestrator/pkg/uploader"
	test_utils "processing-orchestrator/test-utils"
	"testing"
	"time"
)

func TestRecordProcessor_Process_ErrorDuringCooking(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

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
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

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
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

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

// This is a non fatal error
func TestRecordProcessor_Process_ErrorDuringAddToPlaylist(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().AddToPlaylist(mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
		UserInput: processing_common.UserInput{
			Vid: processing_common.VideoOpt{
				PlaylistId: "playlist",
			},
		},
	}
	err := rp.Process(job)
	assert.NoError(t, err)
}

func TestRecordProcessor_Process_ErrorDuringCreatePlaylist(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().CreatePlaylist(mock.Anything).Return(nil, fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
		UserInput: processing_common.UserInput{
			Vid: processing_common.VideoOpt{
				PlaylistId: "",
			},
		},
	}
	err := rp.Process(job)
	assert.NoError(t, err)
}

func TestRecordProcessor_Process_NoPlaylistIdProvided(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 2)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	vid := uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&vid, nil)
	playlist := uploader.Playlist{
		Id: "dd",
	}
	mockUploader.EXPECT().CreatePlaylist(mock.Anything).Return(&playlist, nil)
	mockUploader.EXPECT().AddToPlaylist(vid.Id, playlist.Id).Return(nil)

	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
		UserInput: processing_common.UserInput{
			Vid: processing_common.VideoOpt{
				PlaylistId: "",
			},
		},
	}
	err := rp.Process(job)
	mockUploader.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestRecordProcessor_ProcessWhole(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)
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
	select {
	case evt := <-evtCh:
		pg := evt.ToProgress()
		assert.Equal(t, processing_common.StepDone, pg.Step)
		assert.Equal(t, "prefix", pg.Link)
	case <-time.After(1 * time.Second):
		// Event not received
		t.Fail()
	}
	assert.NoError(t, err)
}

func TestRecordProcessor_ProcessWhole_WithThumbnailAddon(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	mockThumbnailAddon := &test_utils.MockThumbGenService{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{
		ThumbGen: mockThumbnailAddon,
	})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)
	mockThumbnailAddon.EXPECT().GenerateThumbnail(mock.Anything).Return("thumb", nil)
	mockUploader.EXPECT().SetThumbnail(mock.Anything, mock.Anything).Return(nil)
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
	mockStore.AssertExpectations(t)
}

func TestRecordProcessor_ProcessWhole_ErrorThumbnailAddon(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	mockThumbnailAddon := &test_utils.MockThumbGenService{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{
		ThumbGen: mockThumbnailAddon,
	})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)
	mockThumbnailAddon.EXPECT().GenerateThumbnail(mock.Anything).Return("", fmt.Errorf("test"))
	mockUploader.EXPECT().SetThumbnail(mock.Anything, mock.Anything).Return(nil)
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	// An error during thumbnail generation is not fatal
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestRecordProcessor_ProcessWhole_ErrorSetThumbnail(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	mockThumbnailAddon := &test_utils.MockThumbGenService{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{
		ThumbGen: mockThumbnailAddon,
	})

	mockCooker.EXPECT().Cook(mock.Anything, mock.Anything).Return([]string{"cooked"}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockEncoder.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything).Return("video", nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	mockUploader.EXPECT().Upload(mock.Anything, mock.Anything, mock.Anything).Return(&uploader.Video{
		Id:           "fix",
		Title:        "",
		Description:  "",
		CreatedAt:    time.Time{},
		Duration:     0,
		Visibility:   "",
		ThumbnailUrl: "",
		WatchPrefix:  "pre",
	}, nil)

	mockThumbnailAddon.EXPECT().GenerateThumbnail(mock.Anything).Return("thumb", nil)
	mockUploader.EXPECT().SetThumbnail(mock.Anything, mock.Anything).Return(fmt.Errorf("test"))
	job := &job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}
	err := rp.Process(job)
	// An error during thumbnail generation is not fatal
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestRecordProcessor_UpdateInfos(t *testing.T) {
	mockCooker := &test_utils.MockCookingService{}
	mockEncoder := &test_utils.MockEncodingService{}
	mockUploader := &test_utils.MockUploadingService{}
	evtCh := make(chan processing_common.Watchable, 1)
	mockStore := &test_utils.MockJobStore{}
	rp := NewRecordProcessor(mockCooker, mockEncoder, mockUploader, evtCh, mockStore, Addons{})

	mockStore.EXPECT().Get(mock.Anything).Return(&job_store.JobState{
		Id:                 "test",
		Step:               processing_common.StepCooking,
		RawAudioKeys:       []string{"raw"},
		CookedAudioKeys:    nil,
		BackgroundAudioKey: "",
		VideoKey:           "",
	}, nil)
	mockStore.EXPECT().Upsert(mock.Anything).Return(nil)
	err := rp.UpdateInfos("test", processing_common.UserInput{
		Vid: processing_common.VideoOpt{
			Description: "desc",
			Title:       "title",
			Visibility:  processing_common.Public,
		},
	})
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

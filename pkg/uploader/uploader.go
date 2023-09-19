package uploader

import (
	"context"
	"encoding/json"
	"processing-orchestrator/internal/utils"
)

type Uploader struct {
	invoker   utils.Invoker
	component string
}

func NewUploader(invoker utils.Invoker, component string) *Uploader {
	return &Uploader{
		invoker:   invoker,
		component: component,
	}
}

func (u *Uploader) Upload(jobId, storageKey string, opt *VideoOpt) (*Video, error) {
	uploadJob := UploadJob{
		VideoOpt:   *opt,
		StorageKey: storageKey,
		JobId:      jobId,
	}
	bytes, err := json.Marshal(uploadJob)
	if err != nil {
		return nil, err
	}
	res, err := u.invoker.InvokeMethodWithContent(context.Background(), u.component, "v1/videos", "POST", &utils.DataContent{
		ContentType: "application/json",
		Data:        bytes,
	})
	if err != nil {
		return nil, err
	}
	var video Video
	err = json.Unmarshal(res, &video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

package progress_reporter

import (
	"fmt"
	"log/slog"
	processing_common "processing-orchestrator/pkg/processing-common"
	pb "processing-orchestrator/proto"
)

func NewProgressReporter(evtCh chan processing_common.Watchable) *ProgressReporter {
	return &ProgressReporter{
		jobInfos:    map[string][]string{},
		jobWatchers: map[string][]*BidirectionalCom{},
		evtCh:       evtCh,
	}
}

type ProgressReporter struct {
	// Maps id to audioskeys
	jobInfos map[string][]string
	// maps id to watch server
	jobWatchers map[string][]*BidirectionalCom
	evtCh       chan processing_common.Watchable
}

func (pr *ProgressReporter) Start() {
	for {
		evt := <-pr.evtCh
		pg := evt.ToProgress()
		watchers, ok := pr.jobWatchers[pg.JobId]
		if !ok {
			slog.Info(fmt.Sprintf("No watchers for job %s", pg.JobId))
			continue
		}
		status := pr.toStatus(pg)
		for i, channel := range watchers {
			select {
			case <-channel.Canary:
				slog.Debug("Channel closed, skipping")
				pr.Remove(pg.JobId, i)
			default:
				channel.Data <- status
			}
		}
	}
}

func (pr *ProgressReporter) Register(toJobId string, with *BidirectionalCom) {
	pr.jobWatchers[toJobId] = append(pr.jobWatchers[toJobId], with)
}

func (pr *ProgressReporter) AddInfo(jobId string, audiosKeys []string) {
	pr.jobInfos[jobId] = audiosKeys
}

func (pr *ProgressReporter) toStatus(pg *processing_common.ServiceProgress) *pb.ProcessingStatus {

	pb := pb.ProcessingStatus{
		Id:                  pg.JobId,
		Error:               "",
		StepsList:           steps,
		CurrentStepIndex:    uint32(pg.Step),
		ItemList:            nil,
		CurrentItemIndex:    0,
		Done:                false,
		Progress:            pg.Progress,
		Link:                pg.Link,
		CreatedPlaylistLink: pg.CreatedPlaylistLink,
	}

	switch pg.Step {
	case processing_common.StepCooking:
		audiosKeys, ok := pr.jobInfos[pg.JobId]
		if ok {
			pb.ItemList = audiosKeys
			for i, key := range audiosKeys {
				if key == pg.CurrentItem {
					pb.CurrentItemIndex = uint32(i)
					break
				}
			}
		} else {
			slog.Warn(fmt.Sprintf("No audio keys for job %s", pg.JobId))
		}

	case processing_common.StepEncoding:
		pb.ItemList = []string{"videoEnc"}
		pb.CurrentItemIndex = 0
	case processing_common.StepUploading:
		pb.ItemList = []string{"videoUp"}
		pb.CurrentItemIndex = 0
	case processing_common.StepDone:
		pb.Done = true
	}

	// No need
	if pg.Error != nil {
		pb.Error = pg.Error.Error()
		return &pb
	}

	return &pb

}

func (pr *ProgressReporter) Remove(jobId string, index int) {

	watchers := pr.jobWatchers[jobId]
	if pr.jobWatchers[jobId] == nil {
		return
	}
	if index >= 0 && index < len(watchers) {
		watchers = append(watchers[:index], watchers[index+1:]...)
		if len(watchers) == 0 {
			delete(pr.jobWatchers, jobId)
			delete(pr.jobInfos, jobId)
		} else {
			pr.jobWatchers[jobId] = watchers
		}
	}
}

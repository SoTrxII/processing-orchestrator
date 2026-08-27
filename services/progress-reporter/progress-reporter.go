package progress_reporter

import (
	"fmt"
	"log/slog"
	processing_common "processing-orchestrator/pkg/processing-common"
	pb "processing-orchestrator/proto"
	"sync"
)

func NewProgressReporter(evtCh chan processing_common.Watchable) *ProgressReporter {
	return &ProgressReporter{
		jobInfos:    map[string][]string{},
		jobWatchers: map[string][]*BidirectionalCom{},
		evtCh:       evtCh,
	}
}

type ProgressReporter struct {
	// Guards both maps below : clients register from their own gRPC handler
	// while the reporting loop reads and prunes them
	mu sync.Mutex
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
		pr.broadcast(pg.JobId, pr.toStatus(pg))
	}
}

// broadcast hands a status to every client still watching this job, forgetting
// the ones that left. Once a job is over nobody can be interested in it
// anymore, so all of its clients are dropped
func (pr *ProgressReporter) broadcast(jobId string, status *pb.ProcessingStatus) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	watchers, ok := pr.jobWatchers[jobId]
	if !ok {
		slog.Info(fmt.Sprintf("No watchers for job %s", jobId))
		return
	}

	stillWatching := make([]*BidirectionalCom, 0, len(watchers))
	for _, watcher := range watchers {
		if send(watcher, status) {
			stillWatching = append(stillWatching, watcher)
		}
	}

	// The status was handed over already, so a client watching a job that just
	// ended still gets to see it end
	if len(stillWatching) == 0 || status.Done || status.Error != "" {
		pr.forget(jobId)
		return
	}
	pr.jobWatchers[jobId] = stillWatching
}

// send delivers a status to a client and reports whether it is still around to
// receive the next one. A client too slow to keep up loses this update rather
// than holding back every other job
func send(to *BidirectionalCom, status *pb.ProcessingStatus) bool {
	select {
	case <-to.Done:
		slog.Debug("Client stopped watching, dropping it")
		return false
	default:
	}

	select {
	case to.Data <- status:
	default:
		slog.Warn("Client is not reading its progress fast enough, dropping an update")
	}
	return true
}

func (pr *ProgressReporter) Register(toJobId string, with *BidirectionalCom) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	pr.jobWatchers[toJobId] = append(pr.jobWatchers[toJobId], with)
}

func (pr *ProgressReporter) AddInfo(jobId string, audiosKeys []string) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
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
		pr.mu.Lock()
		audiosKeys, ok := pr.jobInfos[pg.JobId]
		pr.mu.Unlock()
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

// forget drops every trace of a job. Callers must hold the lock
func (pr *ProgressReporter) forget(jobId string) {
	delete(pr.jobWatchers, jobId)
	delete(pr.jobInfos, jobId)
}

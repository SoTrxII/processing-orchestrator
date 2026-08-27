// Package demux routes the events published by a remote service to the job
// waiting for them
package demux

import (
	"fmt"
	"log/slog"
	"sync"
)

// Number of events a job can fall behind on before they start being dropped.
// A job reads its events in a tight loop, so a buffer this size only ever
// fills up if the job is gone
const bufferSize = 50

// Demux hands each event to the single job it belongs to.
//
// Several jobs are processed at the same time, and they all share the one
// subscription the service publishes its events on. A channel delivers an
// event to exactly one receiver, so jobs reading a common channel would
// consume each other's events for good : one channel per job and a routing by
// job id is what keeps them apart
type Demux[T any] struct {
	// Guards the map below : jobs register from their own goroutine while
	// events are dispatched from the Dapr topic handler
	mu sync.Mutex
	// One channel per job waiting for its events
	waiters map[string]chan T
}

func New[T any]() *Demux[T] {
	return &Demux[T]{waiters: map[string]chan T{}}
}

// Register opens the channel the events of this job will be delivered on.
// The returned function releases it and must be called once the job is over,
// so the events still coming in for it are dropped rather than piling up
func (d *Demux[T]) Register(jobId string) (<-chan T, func()) {
	events := make(chan T, bufferSize)

	d.mu.Lock()
	defer d.mu.Unlock()
	d.waiters[jobId] = events

	return events, func() {
		d.mu.Lock()
		defer d.mu.Unlock()
		delete(d.waiters, jobId)
	}
}

// Dispatch hands an event to the job it belongs to.
// An event nobody is waiting for anymore (the job errored out, or it finished
// before its last progress event arrived) is dropped rather than blocking the
// Dapr topic handler
func (d *Demux[T]) Dispatch(jobId string, evt T) {
	d.mu.Lock()
	events, isWaiting := d.waiters[jobId]
	d.mu.Unlock()

	if !isWaiting {
		slog.Debug(fmt.Sprintf("[Demux] :: Dropping an event of job %s, nothing is waiting for it", jobId))
		return
	}

	select {
	case events <- evt:
	default:
		slog.Warn(fmt.Sprintf("[Demux] :: Dropping an event of job %s, it isn't reading them fast enough", jobId))
	}
}

// IsWaiting tells whether this job is currently expecting its events
func (d *Demux[T]) IsWaiting(jobId string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, isWaiting := d.waiters[jobId]
	return isWaiting
}

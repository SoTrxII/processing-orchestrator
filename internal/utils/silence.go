package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// resends is how many times a step is asked again of a worker that went
// silent before the job is failed. A worker restarted mid-job loses its
// message for good: the queues are not durable and are acked on delivery
const resends = 2

// Silence fires once a worker has said nothing about a job for `after`
type Silence struct {
	timer  *time.Timer
	after  time.Duration
	resent int
}

func NewSilence(after time.Duration) *Silence {
	return &Silence{timer: time.NewTimer(after), after: after}
}

// C fires once the worker has been silent for too long
func (s *Silence) C() <-chan time.Time {
	return s.timer.C
}

// Heard starts the wait over. Under go 1.21, a timer that fired unread has to
// be drained before Reset, or the stale tick fires right after
func (s *Silence) Heard() {
	if !s.timer.Stop() {
		select {
		case <-s.timer.C:
		default:
		}
	}
	s.timer.Reset(s.after)
}

// AskAgain publishes the job again once the worker went silent, or tells why
// the job is given up
func (s *Silence) AskAgain(worker, jobId string, publish func() error) error {
	if s.resent == resends {
		return fmt.Errorf("%s said nothing about job %s for %s, %d times", worker, jobId, s.after, s.resent+1)
	}
	s.resent++
	slog.Warn(fmt.Sprintf("[Silence] :: no news from %s about job %s for %s, asking again", worker, jobId, s.after))
	if err := publish(); err != nil {
		return err
	}
	s.Heard()
	return nil
}

func (s *Silence) Stop() {
	s.timer.Stop()
}

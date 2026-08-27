package demux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDemux_RoutesEachEventToItsOwnJob(t *testing.T) {
	d := New[string]()
	first, releaseFirst := d.Register("job-1")
	defer releaseFirst()
	second, releaseSecond := d.Register("job-2")
	defer releaseSecond()

	d.Dispatch("job-1", "for the first")
	d.Dispatch("job-2", "for the second")

	assert.Equal(t, "for the first", <-first)
	assert.Equal(t, "for the second", <-second)
}

func TestDemux_NeverHandsAnEventToTheWrongJob(t *testing.T) {
	d := New[string]()
	first, release := d.Register("job-1")
	defer release()

	// A job waiting on a shared channel would consume this one for good
	d.Dispatch("job-2", "for the second")

	assert.Len(t, first, 0)
}

func TestDemux_DropsEventsNobodyIsWaitingFor(t *testing.T) {
	d := New[string]()
	// Must not block nor panic, there is nothing to deliver to
	d.Dispatch("job-1", "for nobody")
	assert.False(t, d.IsWaiting("job-1"))
}

func TestDemux_StopsDeliveringOnceReleased(t *testing.T) {
	d := New[string]()
	events, release := d.Register("job-1")
	release()

	d.Dispatch("job-1", "too late")

	assert.False(t, d.IsWaiting("job-1"))
	assert.Len(t, events, 0)
}

func TestDemux_DropsEventsOfAJobFallingBehind(t *testing.T) {
	d := New[string]()
	events, release := d.Register("job-1")
	defer release()

	// Nothing reads them, so the buffer fills up. The extra events are dropped
	// rather than blocking the topic handler
	for i := 0; i < bufferSize+10; i++ {
		d.Dispatch("job-1", "event")
	}

	assert.Len(t, events, bufferSize)
}

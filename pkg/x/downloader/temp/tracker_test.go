package downloader

import (
	"io"
	"testing"
	"time"

	"github.com/maargenton/go-testpredicate/pkg/verify"
)

func TestTracker(t *testing.T) {
	var tracker = newTracker()
	go func() {
		time.Sleep(time.Millisecond)
		tracker.set(1000)
	}()

	size, err := tracker.get(100)
	verify.That(t, size).Ge(100)
	verify.That(t, err).IsNil()
}

func TestTrackerAdd(t *testing.T) {
	var tracker = newTracker()
	tracker.set(1000)
	go func() {
		time.Sleep(time.Millisecond)
		tracker.add(1000)
	}()

	size, err := tracker.get(1000)
	verify.That(t, size).Ge(200)
	verify.That(t, err).IsNil()
}

func TestTrackerError(t *testing.T) {
	var tracker = newTracker()
	go func() {
		time.Sleep(time.Millisecond)
		tracker.set(10)
		tracker.setError(io.EOF)
	}()

	_, err := tracker.get(100)
	verify.That(t, err).IsError(io.EOF)
}

package proberutil

import (
	"time"

	"testing"

	"github.com/fopina/speedtest-cli/prober"
	"github.com/fopina/speedtest-cli/units"
)

func TestSpeedCollect_NoStream(t *testing.T) {
	grp := prober.NewGroup(1)
	grp.Add(func() (prober.BytesTransferred, error) {
		time.Sleep(10 * time.Millisecond)
		return prober.BytesTransferred(1000), nil
	})

	speed, err := SpeedCollect(grp, nil)
	if err != nil {
		t.Fatalf("SpeedCollect failed: %v", err)
	}
	if speed <= 0 {
		t.Errorf("expected positive speed, got %v", speed)
	}
}

func TestSpeedCollect_WithStream(t *testing.T) {
	grp := prober.NewGroup(1)
	stream := make(chan units.BytesPerSecond, 10)

	grp.Add(func() (prober.BytesTransferred, error) {
		time.Sleep(10 * time.Millisecond)
		return prober.BytesTransferred(1000), nil
	})

	speed, err := SpeedCollect(grp, stream)
	if err != nil {
		t.Fatalf("SpeedCollect failed: %v", err)
	}
	if speed <= 0 {
		t.Errorf("expected positive speed, got %v", speed)
	}

	// Check that some data came through the stream
	select {
	case s := <-stream:
		if s <= 0 {
			t.Errorf("expected positive speed from stream, got %v", s)
		}
	default:
		t.Log("No data in stream yet, but that's ok")
	}

	// Stream should be closed
	_, ok := <-stream
	if ok {
		t.Error("expected stream to be closed")
	}
}

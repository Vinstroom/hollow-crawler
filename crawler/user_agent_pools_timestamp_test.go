package crawler

import (
	"testing"
	"time"
)

func TestNewUserAgentPoolsTimestamp(t *testing.T) {
	ut := NewUserAgentPoolsTimestamp()
	if ut == nil {
		t.Fatal("NewUserAgentPoolsTimestamp() returned nil")
	}
	if ut.timestamps == nil {
		t.Fatal("timestamps map is nil")
	}
}

func TestSetAndGetTimestamp(t *testing.T) {
	ut := NewUserAgentPoolsTimestamp()
	poolID := "test-pool"

	before := time.Now().UnixNano()
	ut.SetTimestamp(poolID)
	after := time.Now().UnixNano()

	ts := ut.GetTimestamp(poolID)
	if ts < before || ts > after {
		t.Errorf("timestamp %d not between %d and %d", ts, before, after)
	}
}

func TestGetTimestampNonExistent(t *testing.T) {
	ut := NewUserAgentPoolsTimestamp()
	ts := ut.GetTimestamp("non-existent")
	if ts != 0 {
		t.Errorf("expected 0 for non-existent pool, got %d", ts)
	}
}

func TestSetTimestampUpdatesExisting(t *testing.T) {
	ut := NewUserAgentPoolsTimestamp()
	poolID := "test-pool"

	ut.SetTimestamp(poolID)
	first := ut.GetTimestamp(poolID)

	time.Sleep(time.Millisecond)
	ut.SetTimestamp(poolID)
	second := ut.GetTimestamp(poolID)

	if second <= first {
		t.Errorf("expected updated timestamp > %d, got %d", first, second)
	}
}

func TestConcurrentAccess(t *testing.T) {
	ut := NewUserAgentPoolsTimestamp()
	done := make(chan struct{})

	go func() {
		for i := 0; i < 100; i++ {
			ut.SetTimestamp("pool-a")
		}
		done <- struct{}{}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			ut.GetTimestamp("pool-b")
		}
		done <- struct{}{}
	}()

	<-done
	<-done
}

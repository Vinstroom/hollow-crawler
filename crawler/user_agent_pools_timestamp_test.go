package crawler

import (
	"sync"
	"testing"
	"time"
)

func TestNewUserAgentPoolsTimestamp(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()
	if uat == nil {
		t.Fatal("NewUserAgentPoolsTimestamp() returned nil")
	}
	if uat.timestamps == nil {
		t.Fatal("timestamps map is nil")
	}
}

func TestSetAndGetTimestamp(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()
	poolID := "test-pool"

	before := time.Now().UnixNano()
	uat.SetTimestamp(poolID)
	after := time.Now().UnixNano()

	ts := uat.GetTimestamp(poolID)
	if ts < before || ts > after {
		t.Errorf("timestamp %d should be between %d and %d", ts, before, after)
	}
}

func TestGetTimestampDefaultZero(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()
	ts := uat.GetTimestamp("nonexistent")
	if ts != 0 {
		t.Errorf("expected 0 for nonexistent pool, got %d", ts)
	}
}

func TestGetTimestampAfterSet(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()
	poolID := "pool-1"

	uat.SetTimestamp(poolID)
	ts1 := uat.GetTimestamp(poolID)
	time.Sleep(10 * time.Millisecond)
	uat.SetTimestamp(poolID)
	ts2 := uat.GetTimestamp(poolID)

	if ts2 <= ts1 {
		t.Errorf("second timestamp %d should be greater than first %d", ts2, ts1)
	}
}

func TestConcurrentAccess(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			poolID := "pool"
			uat.SetTimestamp(poolID)
			uat.GetTimestamp(poolID)
		}(i)
	}
	wg.Wait()
}

func TestMultiplePools(t *testing.T) {
	uat := NewUserAgentPoolsTimestamp()

	uat.SetTimestamp("pool-a")
	uat.SetTimestamp("pool-b")
	tsA := uat.GetTimestamp("pool-a")
	tsB := uat.GetTimestamp("pool-b")

	if tsA == 0 || tsB == 0 {
		t.Fatal("timestamps should be non-zero")
	}

	zeroTS := uat.GetTimestamp("pool-c")
	if zeroTS != 0 {
		t.Errorf("expected 0 for unset pool-c, got %d", zeroTS)
	}
}

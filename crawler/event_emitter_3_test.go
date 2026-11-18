package crawler

import (
	"context"
	"testing"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestEventEmitter(t *testing.T) {
	ctx := context.Background()
	e := NewEventEmitter(ctx)

	var emitEvents []Event
	emittedMutex := sync.Mutex{}
	e.SetEmitFunc(func(event Event) {
		emittedMutex.Lock()
		defer emittedMutex.Unlock()
		emitEvents = append(emitEvents, event)
	})

	e.EmitDeduplicate("TestDeduplicateData")
	e.EmitExport("TestExportData")

	e.Close()

	emittedMutex.Lock()
	defer emittedMutex.Unlock()

	// The events emitted before the cancellation should be executed
	assert.Len(t, emitEvents, 2)
	assert.Equal(t, EventDeduplicate, emitEvents[0].EventType)
	assert.Equal(t, "TestDeduplicateData", emitEvents[0].Data)
	assert.Equal(t, EventExport, emitEvents[1].EventType)
	assert.Equal(t, "TestExportData", emitEvents[1].Data)

	// No more events are expected after cancellation
	e.EmitDeduplicate("AdditionalDeduplicateData")
	emittedMutex.Lock()
	defer emittedMutex.Unlock()
	assert.Len(t, emitEvents, 2) // Still only 2 events

	// After Close() is called, calling SetEmitFunc again should have no effect
	e.SetEmitFunc(func(Event) {})
	e.EmitDeduplicate("NoEffectData")
	emittedMutex.Lock()
	defer emittedMutex.Unlock()
	assert.Len(t, emitEvents, 2) // Still only 2 events
}

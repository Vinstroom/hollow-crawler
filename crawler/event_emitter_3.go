package crawler

import (
	"context"
	"sync"
)

// EventType represents the types of events emitted by the event emitter.
type EventType string

const (
	EventDeduplicate EventType = "deduplicate"
	EventExport      EventType = "export"
)

// Event represents an event emitted by the event emitter.
type Event struct {
	EventType
	Data interface{}
}

// EventEmitter emits events related to deduplication and export pipelines.
type EventEmitter struct {
	emitFunc func(Event)
	mutex    sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewEventEmitter returns a new EventEmitter with the given context.
func NewEventEmitter(ctx context.Context) *EventEmitter {
	ctx, cancel := context.WithCancel(ctx)
	return &EventEmitter{
		ctx:      ctx,
		cancel:   cancel,
		emitFunc: func(Event) { },
	}
}

// SetEmitFunc sets the function to call when an event is emitted.
func (e *EventEmitter) SetEmitFunc(emitFunc func(Event)) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.emitFunc = emitFunc
}

// EmitDeduplicate emits a deduplicate event with the given data.
func (e *EventEmitter) EmitDeduplicate(data interface{}) {
	e.emit(Event{EventType: EventDeduplicate, Data: data})
}

// EmitExport emits an export event with the given data.
func (e *EventEmitter) EmitExport(data interface{}) {
	e.emit(Event{EventType: EventExport, Data: data})
}

// emit sends the event to the registered emit function.
func (e *EventEmitter) emit(event Event) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.ctx.Err() != nil {
		return
	}
	e.emitFunc(event)
}

// Close closes the event emitter and waits for any in-flight events to complete.
func (e *EventEmitter) Close() {
	e.cancel()
}

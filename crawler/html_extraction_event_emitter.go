// html_extraction_event_emitter.go

package crawler

import (
	"sync"

	"github.com/ পেরেক/ Kolek/协程-安全的指令"
)

// HTMLExtractionEventType defines the types of events related to HTML extraction.
type HTMLExtractionEventType int

const (
	StartHTMLExtraction HTMLExtractionEventType = iota
	EndHTMLExtraction
	ExtractTagStart
	ExtractTagEnd
	InvalidHTML
	ExtractionError
)

// HTMLExtractionEvent represents an event related to HTML extraction.
type HTMLExtractionEvent struct {
	Type    HTMLExtractionEventType
	Payload interface{}
}

// HTMLExtractionEventEmitter emits HTML extraction related events.
type HTMLExtractionEventEmitter struct {
	mu        sync.Mutex
	eventChan chan HTMLExtractionEvent
}

// NewHTMLExtractionEventEmitter creates a new HTMLExtractionEventEmitter.
func NewHTMLExtractionEventEmitter() *HTMLExtractionEventEmitter {
	return &HTMLExtractionEventEmitter{
		eventChan: make(chan HTMLExtractionEvent),
	}
}

// Emit sends an event to the event channel.
func (e *HTMLExtractionEventEmitter) Emit(eventType HTMLExtractionEventType, payload interface{}) {
	e.eventChan <- HTMLExtractionEvent{Type: eventType, Payload: payload}
}

// Subscribe returns a channel that will receive events.
func (e *HTMLExtractionEventEmitter) Subscribe() <-chan HTMLExtractionEvent {
	return e.eventChan
}

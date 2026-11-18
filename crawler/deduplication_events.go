package crawler

// DeduplicationEventType represents the type of deduplication events.
type DeduplicationEventType int

const (
	// DeduplicationEventTypeAdded indicates that a URL was added to the deduplication state.
	DeduplicationEventTypeAdded DeduplicationEventType = iota

	// DeduplicationEventTypeRemoved indicates that a URL was removed from the deduplication state.
	DeduplicationEventTypeRemoved

	// DeduplicationEventTypeIgnored indicates that a URL was ignored by the deduplication process.
	DeduplicationEventTypeIgnored
)

// DeduplicationEvent represents an event related to the deduplication process.
type DeduplicationEvent struct {
	eventType DeduplicationEventType
	url        string
}

// DeduplicationEvents is a channel that emits DeduplicationEvents.
type DeduplicationEvents chan DeduplicationEvent

// Emit emits a DeduplicationEvent onto the channel.
func (e DeduplicationEvents) Emit(eventType DeduplicationEventType, url string) {
	e <- DeduplicationEvent{eventType, url}
}

// DeduplicationPipelineEventType represents the type of deduplication pipeline events.
type DeduplicationPipelineEventType int

const (
	// DeduplicationPipelineEventTypeStarted indicates that the deduplication pipeline has started.
	DeduplicationPipelineEventTypeStarted DeduplicationPipelineEventType = iota

	// DeduplicationPipelineEventTypeFinished indicates that the deduplication pipeline has finished.
	DeduplicationPipelineEventTypeFinished

	// DeduplicationPipelineEventTypeFailed indicates that the deduplication pipeline has failed.
	DeduplicationPipelineEventTypeFailed
)

// DeduplicationPipelineEvent represents an event related to the deduplication pipeline.
type DeduplicationPipelineEvent struct {
	eventType DeduplicationPipelineEventType
	err       error
}

// DeduplicationPipelineEvents is a channel that emits DeduplicationPipelineEvents.
type DeduplicationPipelineEvents chan DeduplicationPipelineEvent

// Emit emits a DeduplicationPipelineEvent onto the channel.
func (e DeduplicationPipelineEvents) Emit(eventType DeduplicationPipelineEventType, err error) {
	e <- DeduplicationPipelineEvent{eventType, err}
}

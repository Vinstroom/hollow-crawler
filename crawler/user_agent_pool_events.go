package crawler

// UserAgentPoolEvent is an event emitted by the user-agent pool.
type UserAgentPoolEvent struct {
	EventType  string // e.g., "AgentUsed", "AgentPooled"
	UserAgent  string
	Pool       string // name of the pool
	Timestamp  time.Time
	Details    string // optional details about the event
}

// UserAgentPoolEvents is a channel that emits UserAgentPoolEvent events.
type UserAgentPoolEvents chan UserAgentPoolEvent

// EmitEvent sends an event to the channel. If the channel is closed, the event
// is discarded.
func (e UserAgentPoolEvents) EmitEvent(eventType, userAgent, pool string, timestamp time.Time, details string) {
	select {
	case e <- UserAgentPoolEvent{EventType: eventType, UserAgent: userAgent, Pool: pool, Timestamp: timestamp, Details: details}:
	default:
		// Ignore if channel is closed
	}
}

// Close closes the channel, preventing any further events from being sent.
func (e UserAgentPoolEvents) Close() {
	close(e)
}

// UserAgentPoolEventListener is an interface that can receive user-agent pool events.
type UserAgentPoolEventListener interface {
	OnUserAgentPoolEvent(event UserAgentPoolEvent)
}

// userAgentPoolEventEmitter emits user-agent pool events to listeners.
type userAgentPoolEventEmitter struct {
	events UserAgentPoolEvents
	listeners []UserAgentPoolEventListener
}

// newUserAgentPoolEventEmitter creates a new user-agent pool event emitter.
func newUserAgentPoolEventEmitter() *userAgentPoolEventEmitter {
	return &userAgentPoolEventEmitter{
		events: make(UserAgentPoolEvents),
	}
}

// Start starts emitting events. Returns a channel that will receive events.
func (em *userAgentPoolEventEmitter) Start() <-chan UserAgentPoolEvent {
	return em.events
}

// AddListener adds a listener that will receive events.
func (em *userAgentPoolEventEmitter) AddListener(listener UserAgentPoolEventListener) {
	em.listeners = append(em.listeners, listener)
}

// EmitEvent sends an event to all listeners.
func (em *userAgentPoolEventEmitter) EmitEvent(eventType, userAgent, pool string, timestamp time.Time, details string) {
	em.events.EmitEvent(eventType, userAgent, pool, timestamp, details)
	for _, listener := range em.listeners {
		listener.OnUserAgentPoolEvent(UserAgentPoolEvent{EventType: eventType, UserAgent: userAgent, Pool: pool, Timestamp: timestamp, Details: details})
	}
}

// Close closes the event emitter, preventing any further events from being sent.
func (em *userAgentPoolEventEmitter) Close() {
	em.events.Close()
}

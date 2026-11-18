package crawler

// RobotsHandlerEvents represents events that can be emitted during robots handling.
type RobotsHandlerEvents struct {
	EventEmitter
}

// EmitRobotPolicyFound emits when a robot policy is found.
func (e *RobotsHandlerEvents) EmitRobotPolicyFound(url string) {
	e.emit(Event{
		Type: "robotpolicyfound",
		Data: map[string]interface{}{
			"url": url,
		},
	})
}

// EmitRobotPolicyNotFound emits when no robot policy is found.
func (e *RobotsHandlerEvents) EmitRobotPolicyNotFound(url string) {
	e.emit(Event{
		Type: "robotpolicynotfound",
		Data: map[string]interface{}{
			"url": url,
		},
	})
}

// EmitRobotPolicyViolation emits when a robot policy is violated.
func (e *RobotsHandlerEvents) EmitRobotPolicyViolation(url string, policy RobotPolicy, violation string) {
	e.emit(Event{
		Type: "robotpolicypolicyviolation",
		Data: map[string]interface{}{
			"url":      url,
			"policy":   policy.String(),
			"violation": violation,
		},
	})
}

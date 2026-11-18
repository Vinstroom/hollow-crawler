package crawler

import (
	"sync"
	"time"
)

const (
	UAIdle State = iota
	UAActive
	UAExhausted
)

type State int

type ProxyRotationSM struct {
	state       State
	proxyCount  int
	proxyTick   time.Duration
	currentProxy Proxy
	userAgentPool UserAgentPoolInterface
	mtx         sync.Mutex
}

type ProxyRotatorInterface interface {
	NextProxy() Proxy
}

type UserAgentPoolInterface interface {
	NextUserAgent() string
}

func NewProxyRotationSM(initialProxy Proxy, count int, tick time.Duration, userAgentPool UserAgentPoolInterface) *ProxyRotationSM {
	return &ProxyRotationSM{
		state:        UAIdle,
		proxyCount:   count,
		proxyTick:    tick,
		currentProxy: initialProxy,
		userAgentPool: userAgentPool,
	}
}

func (s *ProxyRotationSM) SetState(state State) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.state = state
}

func (s *ProxyRotationSM) GetState() State {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.state
}

func (s *ProxyRotationSM) SetProxy(proxy Proxy) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.currentProxy = proxy
}

func (s *ProxyRotationSM) GetProxy() Proxy {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.currentProxy
}

func (s *ProxyRotationSM) Rotate(proxyRotator ProxyRotatorInterface) {
	switch s.state {
	case UAIdle:
		s.currentProxy = proxyRotator.NextProxy()
		s.currentProxy.UserAgent = s.userAgentPool.NextUserAgent()
		s.state = UAActive
		time.AfterFunc(s.proxyTick, func() {
			s.Rotate(proxyRotator)
		})
	case UAActive:
		// Continue rotating every 'tick' duration
		s.currentProxy = proxyRotator.NextProxy()
		s.currentProxy.UserAgent = s.userAgentPool.NextUserAgent()
		time.AfterFunc(s.proxyTick, func() {
			s.Rotate(proxyRotator)
		})
	case UAExhausted:
		// Optionally start rotating again after a delay
		time.AfterFunc(5*time.Minute, func() {
			s.state = UAIdle
			s.Rotate(proxyRotator)
		})
	}
}

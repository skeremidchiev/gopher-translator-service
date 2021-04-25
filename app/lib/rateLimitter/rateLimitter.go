// Simple rate limitter for external API calls in order to prevent bursts and requests drops
package rateLimitter

import (
	"sync"
	"time"
)

type rateLimitter struct {
	// number of calls in current interval
	callsTracker int
	// calls limit
	maxCalls int
	// timeout when maxCalls is reached
	timeout time.Duration
	// mutex protects callsTracker variable
	m sync.Mutex
}

// NewRateLimitter
// maxRequests -> number of maximum simultaneous requests
// timeOutDuration -> wait time after limit is reached
func NewRateLimitter(maxRequests int, timeOutDuration time.Duration) *rateLimitter {
	return &rateLimitter{
		callsTracker: 0,
		timeout:      timeOutDuration,
		maxCalls:     maxRequests,
	}
}

// Throttle
func (rl *rateLimitter) Throttle() {
	rl.m.Lock()
	defer rl.m.Unlock()

	rl.callsTracker++
	if rl.callsTracker > rl.maxCalls {
		<-time.After(rl.timeout)
		rl.callsTracker = 0
	}
}

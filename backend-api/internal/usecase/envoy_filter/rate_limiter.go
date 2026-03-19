package envoy_filter

import (
	"sync"
	"time"
)

type rateCounter struct {
	windowStart time.Time
	count       int32
}

type InMemoryRateLimiter struct {
	mu       sync.Mutex
	counters map[string]*rateCounter
}

func NewInMemoryRateLimiter() *InMemoryRateLimiter {
	return &InMemoryRateLimiter{counters: make(map[string]*rateCounter)}
}

func (r *InMemoryRateLimiter) Allow(key string, limit int32, period time.Duration, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter, ok := r.counters[key]
	if !ok {
		r.counters[key] = &rateCounter{windowStart: now, count: 1}
		return true
	}

	if now.Sub(counter.windowStart) >= period {
		counter.windowStart = now
		counter.count = 1
		return true
	}

	if counter.count >= limit {
		return false
	}

	counter.count++
	return true
}

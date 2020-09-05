package ratelimter

import (
	"context"
	"sync/atomic"
	"time"
)

type RateLimiter interface {
	Start()

	TryAcquire() bool

	Acquire() bool

	Stop()
}

type StableRateLimiter struct {
	threshold        int64
	currentThreshold int64
	refillPeriod     time.Duration
	broadcastChannel chan bool
	exitChannel      chan bool
}

func NewRateLimiter(threshold int64, refillPeriod time.Duration) (rateLimiter *StableRateLimiter) {
	return &StableRateLimiter{
		threshold:        threshold,
		currentThreshold: threshold,
		refillPeriod:     refillPeriod,
		broadcastChannel: make(chan bool),
		exitChannel:      make(chan bool),
	}
}

func (l *StableRateLimiter) Start() {
	ec := l.exitChannel
	go func() {
		for true {
			select {
			case <-ec:
				return
			default:
				atomic.StoreInt64(&l.currentThreshold, l.threshold)
				time.Sleep(l.refillPeriod)
				close(l.broadcastChannel)
				l.broadcastChannel = make(chan bool)
			}
		}
	}()
}

func (l *StableRateLimiter) TryAcquire() (allow bool) {
	permit := atomic.AddInt64(&l.currentThreshold, -1)
	if permit < 0 {
		allow = false
	} else {
		allow = true
	}
	return allow
}

func (l *StableRateLimiter) Acquire() (allow bool) {
	for {
		permit := atomic.AddInt64(&l.currentThreshold, -1)
		if permit > 0 {
			allow = false
			select {
			case <-l.broadcastChannel:
				continue
			}
		}
		return true
	}
}

func (l *StableRateLimiter) Stop() {
	close(l.exitChannel)
}

func (l *StableRateLimiter) AcquireContext(ctx context.Context) (allow bool) {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}

		permit := atomic.AddInt64(&l.currentThreshold, -1)
		if permit < 0 {
			allow = false
			select {
			case <-ctx.Done():
				return false
			case <-l.broadcastChannel:
				continue
			}
		}
		return true
	}
}

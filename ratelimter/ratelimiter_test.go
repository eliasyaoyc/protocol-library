package ratelimter

import (
	"context"
	"testing"
	"time"
)

func TestStableRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(1, 10*time.Millisecond)
	limiter.Start()
	defer limiter.Stop()

	acquire := limiter.TryAcquire()
	if !acquire {
		t.Error("Unexpected blocked by rate limier")
	}
	acquire = limiter.TryAcquire()
	if acquire {
		t.Error("Should be blocked")
	}
}

func TestNoBlockMode(t *testing.T) {
	limiter := NewRateLimiter(5, 1*time.Second)
	for {
		allow := limiter.TryAcquire()
		if allow {
			// do something
		}
	}
}

func TestBlockWithContext(t *testing.T) {
	for {
		limiter := NewRateLimiter(10, 1*time.Second)
		for {
			allow := limiter.AcquireContext(context.Background())
			if allow {
				// do something
			}
		}
	}
}

func TestBlockMode(t *testing.T) {
	limiter := NewRateLimiter(10, 1*time.Second)
	for {
		allow := limiter.Acquire()
		if allow {
			// do something
		}
	}
}

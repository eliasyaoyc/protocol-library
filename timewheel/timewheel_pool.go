package timewheel

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type TimeWheelPool struct {
	pool []*TimeWheel
	size int64
	incr int64
}

func NewTimeWheelPool(size int, tick time.Duration, bucketsNum int, options ...optionCall) (*TimeWheelPool, error) {
	twp := &TimeWheelPool{
		pool: make([]*TimeWheel, size),
		size: int64(size),
	}
	for i := 0; i < bucketsNum; i++ {
		tw, err := NewTimeWheel(tick, bucketsNum, options...)
		if err != nil {
			return twp, err
		}
		twp.pool[i] = tw
	}
	return twp, nil
}

func (twp *TimeWheelPool) Get() *TimeWheel {
	incr := atomic.AddInt64(&twp.incr, 1)
	idx := incr % twp.size
	return twp.pool[idx]
}

func (twp *TimeWheelPool) GetRandom() *TimeWheel {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Int63n(twp.size)
	return twp.pool[idx]
}

func (twp *TimeWheelPool) Start() {
	for _, tw := range twp.pool {
		tw.Start()
	}
}

func (twp TimeWheelPool) Stop() {
	for _, tw := range twp.pool {
		tw.Stop()
	}
}

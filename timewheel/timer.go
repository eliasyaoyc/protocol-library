package timewheel

import (
	"errors"
	"sync"
	"time"
)

const (
	typeTimer taskType = iota
	typeTicker

	modeIsCircle  = true
	modeNotCircle = false

	modeIsAsync  = true
	modeNotAsync = false
)

type taskType int64
type taskID int64

type Task struct {
	delay    time.Duration
	id       taskID
	round    int
	callback func()

	async  bool
	stop   bool
	circle bool
	// circleNum int
}

// for sync.Pool
func (t *Task) Reset() {
	t.round = 0
	t.callback = nil

	t.async = false
	t.stop = false
	t.circle = false
}

type optionCall func(*TimeWheel) error

func TickSafeMode() optionCall {
	return func(o *TimeWheel) error {
		o.tickQueue = make(chan time.Time, 10)
		return nil
	}
}

func SetSyncPool(state bool) optionCall {
	return func(o *TimeWheel) error {
		o.syncPool = state
		return nil
	}
}

type TimeWheel struct {
	randomID int64

	tick      time.Duration
	ticker    *time.Ticker
	tickQueue chan time.Time

	bucketsNum    int
	buckets       []map[taskID]*Task // key: added item, value: *Task
	bucketIndexes map[taskID]int     // key: added item, value: bucket position

	currentIndex int

	onceStart sync.Once

	addC    chan *Task
	removeC chan *Task
	stopC   chan struct{}

	exited   bool
	syncPool bool
}

// NewTimeWheel create new time wheel
func NewTimeWheel(tick time.Duration, bucketsNum int, options ...optionCall) (*TimeWheel, error) {
	if tick.Seconds() < 0.1 {
		return nil, errors.New("invalid params, must tick >= 100 ms")
	}
	if bucketsNum <= 0 {
		return nil, errors.New("invalid params, must bucketsNum > 0")
	}

	tw := &TimeWheel{
		// tick
		tick:      tick,
		tickQueue: make(chan time.Time, 10),

		// store
		bucketsNum:    bucketsNum,
		bucketIndexes: make(map[taskID]int, 1024*100),
		buckets:       make([]map[taskID]*Task, bucketsNum),
		currentIndex:  0,

		// signal
		addC:    make(chan *Task, 1024*5),
		removeC: make(chan *Task, 1024*2),
		stopC:   make(chan struct{}),
	}

	for i := 0; i < bucketsNum; i++ {
		tw.buckets[i] = make(map[taskID]*Task, 16)
	}

	for _, op := range options {
		op(tw)
	}

	return tw, nil
}

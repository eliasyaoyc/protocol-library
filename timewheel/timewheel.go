package timewheel

import "time"

var (
	DefaultTimewWheel, _ = NewTimeWheel(time.Second, 120)
)

func init() {
	DefaultTimewWheel.Start()
}

func ResetDefaultTimeWheel(tw *TimeWheel) {
	tw.Start()
	DefaultTimewWheel = tw
}

func Add(delay time.Duration, callback func()) *Task {
	return DefaultTimewWheel.Add(delay, callback)
}

func AddCron(delay time.Duration, callback func()) *Task {
	return DefaultTimewWheel.AddCron(delay, callback)
}

func Remove(task *Task) {
	DefaultTimewWheel.remove(task)
}

func NewTimer(delay time.Duration) *Timer {
	return DefaultTimewWheel.NewTimer(delay)
}

func NewTicker(delay time.Duration) *Ticker {
	return DefaultTimewWheel.NewTicker(delay)
}

func AfterFunc(delay time.Duration, callback func()) *Timer {
	return DefaultTimewWheel.AfterFunc(delay, callback)
}

func After(delay time.Duration) <-chan time.Time {
	return DefaultTimewWheel.After(delay)
}

func Sleep(delay time.Duration) {
	DefaultTimewWheel.Sleep(delay)
}

package timewheel

import "time"

var (
	DefaultTimewWheel,_ = NewTimeWheel(time.Second,120)
)



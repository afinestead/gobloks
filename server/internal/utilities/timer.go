package utilities

import (
	"time"
)

type Timer struct {
	timer    *time.Timer
	duration time.Duration
	start    time.Time
	elapsed  time.Duration
	callback func()
}

func InitTimer(d time.Duration, callback func()) *Timer {
	return &Timer{
		timer:    nil,
		duration: d,
		start:    time.Unix(0, 0),
		elapsed:  0,
		callback: callback,
	}
}

func (t *Timer) Cancel() bool {
	return t.timer.Stop()
}

func (t *Timer) Pause() {
	if !t.timer.Stop() {
		<-t.timer.C
	}
	t.elapsed = time.Since(t.start)
}

func (t *Timer) Start() {
	t.duration -= t.elapsed
	t.start = time.Now()
	t.timer = time.NewTimer(t.duration * time.Second)

	go func() {
		<-t.timer.C // Wait for the timer to expire
		if t.callback != nil {
			t.callback()
		}
	}()
}

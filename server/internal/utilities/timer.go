package utilities

import (
	"fmt"
	"sync"
	"time"
)

type Timer struct {
	timer        *time.Timer
	remaining    time.Duration
	start        time.Time
	bonus        time.Duration
	callback     func(args ...any)
	callbackArgs []any
	expired      bool
	mtx          sync.Mutex
}

func InitTimer(ms, bonus uint, callback func(args ...any), args ...any) *Timer {
	return &Timer{
		timer:        nil,
		remaining:    time.Duration(ms) * time.Millisecond,
		start:        time.Unix(0, 0),
		bonus:        time.Duration(bonus) * time.Millisecond,
		callback:     callback,
		callbackArgs: args,
		expired:      false,
		mtx:          sync.Mutex{},
	}
}

func (t *Timer) Cancel() bool {
	return t.timer.Stop()
}

func (t *Timer) Pause() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.timer == nil { // for initial start
		return
	}

	if !t.timer.Stop() {
		<-t.timer.C
	}
	// Add bonus time when pausing
	elapsed := (time.Since(t.start) - t.bonus)
	t.remaining -= elapsed
}

func (t *Timer) Start() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	fmt.Println(t.remaining)

	t.start = time.Now()
	t.timer = time.NewTimer(t.remaining)

	go func() {
		<-t.timer.C // Wait for the timer to expire
		t.mtx.Lock()
		t.remaining = 0
		t.expired = true
		t.mtx.Unlock()
		if t.callback != nil {
			t.callback(t.callbackArgs...)
		}
	}()
}

func (t *Timer) TimeLeftMs() uint {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return uint(t.remaining / time.Millisecond)
}

func (t *Timer) Expired() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.expired
}

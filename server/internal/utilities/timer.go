package utilities

import (
	"sync"
	"time"
)

type Timer struct {
	timer        *time.Timer
	remaining    time.Duration
	start        time.Time
	bonus        time.Duration
	expired      bool
	callback     func(args ...any)
	callbackArgs []any
	mtx          sync.Mutex
}

func InitTimer(ms, bonus uint, callback func(args ...any), args ...any) *Timer {
	return &Timer{
		timer:        nil,
		remaining:    time.Duration(ms) * time.Millisecond,
		start:        time.Unix(0, 0),
		bonus:        time.Duration(bonus) * time.Millisecond,
		expired:      false,
		callback:     callback,
		callbackArgs: args,
		mtx:          sync.Mutex{},
	}
}

func (t *Timer) Cancel() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	ok := true
	if t.timer != nil {
		ok = t.timer.Stop()
		t.timer = nil
	}
	return ok
}

func (t *Timer) Pause() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.timer == nil {
		return
	}

	if !t.timer.Stop() {
		<-t.timer.C
	}
	// Add bonus time when pausing
	elapsed := (time.Since(t.start) - t.bonus)
	t.remaining -= elapsed
	t.timer = nil
}

func (t *Timer) Start() {
	t.mtx.Lock()
	defer t.mtx.Unlock()

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
	if t.expired {
		return 0
	} else if t.timer == nil {
		return uint(t.remaining / time.Millisecond)
	} else {
		return uint((t.remaining - time.Since(t.start)) / time.Millisecond)
	}
}

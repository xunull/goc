package guard_pool

import (
	"sync/atomic"
	"time"
)

type TempoLimiter struct {
	Limit      int32
	SecondSpan int32
	startTime  time.Time
	count      *int32
	BaseLimiter
}

func NewTempoLimit(limit, second int32) *TempoLimiter {
	tl := &TempoLimiter{
		Limit:      limit,
		SecondSpan: second,
	}
	var t int32
	t = 0
	tl.count = &t
	go func() {
		tl.start()
	}()
	return tl
}

func (t *TempoLimiter) start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		return
	}
	t.running = true
	ticker := time.NewTicker(time.Duration(t.SecondSpan) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				atomic.SwapInt32(t.count, 0)
				t.mu.Lock()
				t.restrict = false
				go func() {
					t.callSubscribeForResetFunc()
				}()
				t.mu.Unlock()

			}
		}
	}()
}

func (t *TempoLimiter) subscribeForReset(f func()) {
	t.subscribeForResetFunc = f
}

func (t *TempoLimiter) subscribeForRestrict(f func()) {
	t.subscribeForRestrictFunc = f
}

func (t *TempoLimiter) Hit() bool {
	if t.restrict {
		return false
	}
	n := atomic.AddInt32(t.count, 1)
	if n > t.Limit {
		t.mu.Lock()
		t.restrict = true
		go func() {
			t.callSubscribeForRestrictFunc()
		}()
		t.mu.Unlock()
		return false
	} else {
		return true
	}
}

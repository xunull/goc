package guard_pool

import (
	"sync/atomic"
	"time"
)

type IntervalLimiter struct {
	interval time.Duration
	last     time.Time
	token    int32
	BaseLimiter
}

func NewIntervalLimiter(interval time.Duration) *IntervalLimiter {
	il := &IntervalLimiter{
		interval: interval,
	}

	go func() {
		il.start()
	}()
	return il
}

func (i *IntervalLimiter) hit() bool {

	swap := atomic.CompareAndSwapInt32(&i.token, 0, 1)
	if swap {
		return true
	} else {

		go func() {
			i.callSubscribeForRestrictFunc()
		}()
		return false
	}
}

func (i *IntervalLimiter) start() {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.running {
		return
	}
	i.running = true
	ticker := time.NewTicker(i.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				swap := atomic.CompareAndSwapInt32(&i.token, 1, 0)
				if swap {
					i.callSubscribeForResetFunc()
				}
			}
		}
	}()
}

func (i *IntervalLimiter) subscribeForRestrict(f func()) {
	i.subscribeForRestrictFunc = f
}

func (i *IntervalLimiter) subscribeForReset(f func()) {
	i.subscribeForResetFunc = f
}

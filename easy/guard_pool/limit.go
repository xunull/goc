package guard_pool

import "sync"

type BaseLimiter struct {
	mu                       sync.Mutex
	running                  bool
	restrict                 bool
	subscribeForResetFunc    func()
	subscribeForRestrictFunc func()
}

func (b *BaseLimiter) callSubscribeForResetFunc() {
	go func() {
		if b.subscribeForResetFunc != nil {
			b.subscribeForResetFunc()
		}
	}()
}

func (b *BaseLimiter) callSubscribeForRestrictFunc() {

	go func() {
		if b.subscribeForRestrictFunc != nil {
			b.subscribeForRestrictFunc()
		}
	}()
}

type limiter interface {
	hit() bool
	start()
	subscribeForRestrict(f func())
	subscribeForReset(f func())
}

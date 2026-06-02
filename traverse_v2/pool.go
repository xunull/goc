package traverse_v2

import "sync"

// pool is a fixed-worker goroutine pool with a bounded queue.
//
// submit is non-blocking: if the queue is full it spawns an overflow
// goroutine instead. This is deliberate. The dir pool's tasks recursively
// submit child tasks; if submit blocked, a saturated pool could deadlock
// when all workers are simultaneously stuck inside submit (waiting on a
// queue that only they can drain). Overflow trades a temporary spike in
// goroutine count for a guaranteed forward-progress property.
type pool struct {
	queue   chan func()
	closeMu sync.Once
}

func newPool(workers, queueSize int) *pool {
	if workers < 1 {
		workers = 1
	}
	if queueSize < workers {
		queueSize = workers
	}
	p := &pool{queue: make(chan func(), queueSize)}
	for i := 0; i < workers; i++ {
		go p.run()
	}
	return p
}

func (p *pool) run() {
	for task := range p.queue {
		task()
	}
}

func (p *pool) submit(task func()) {
	select {
	case p.queue <- task:
	default:
		go task()
	}
}

func (p *pool) close() {
	p.closeMu.Do(func() { close(p.queue) })
}

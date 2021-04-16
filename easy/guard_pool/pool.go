package guard_pool

import "github.com/xunull/goc/easy"

type TaskFunc func()

const (
	Running = iota
	Pending
	Over
	Idle
)

type GoWorker struct {
	id     string
	task   TaskFunc
	status int
}

type GuardPool struct {
	count        int
	RunningSheet map[string]*GoWorker
	IdleSheet    map[string]*GoWorker
	RunningChan  chan *GoWorker
	IdleChan     chan *GoWorker
	StatusChan   chan *GoWorker
	TaskFuncChan chan TaskFunc
	idMarker     *easy.IdMarker
	pauseChan    chan struct{}
	running      bool
	runningCount int
}

func NewPool(count int) *GuardPool {
	pool := &GuardPool{
		RunningSheet: make(map[string]*GoWorker, count),
		IdleSheet:    make(map[string]*GoWorker, count),
		RunningChan:  make(chan *GoWorker, count),
		IdleChan:     make(chan *GoWorker, count),
		StatusChan:   make(chan *GoWorker, count),
		TaskFuncChan: make(chan TaskFunc, count),
		count:        count,
		idMarker:     &easy.IdMarker{},
	}

	pool.prepareWorker()
	return pool
}

func (r *GuardPool) prepareWorker() {
	for i := 0; i < r.count; i++ {
		w := &GoWorker{
			id: r.idMarker.GetNewWorkerId(),
		}
		r.IdleChan <- w
	}
}

func (r *GuardPool) GetRunningCount() int {
	return r.runningCount
}

func (r *GuardPool) runStatusChan() {
	go func() {
		for w := range r.StatusChan {
			if w.status == Pending {
				delete(r.IdleSheet, w.id)
				w.status = Running
				r.RunningSheet[w.id] = w
				r.runningCount += 1
				r.RunningChan <- w
			} else if w.status == Over {
				delete(r.RunningSheet, w.id)
				r.runningCount -= 1
				w.status = Idle
				r.IdleSheet[w.id] = w
				r.IdleChan <- w
			}
		}
	}()
}

func (r *GuardPool) runWorkCore() {

	for w := range r.RunningChan {
		if w.status == Pending {
			r.StatusChan <- w
		} else {
			go func(t *GoWorker) {
				t.task()
				t.task = nil
				t.status = Over
				r.StatusChan <- t
			}(w)
		}
	}
}

func (r *GuardPool) runTaskCore() {
	for {
		select {
		case task := <-r.TaskFuncChan:
			w := <-r.IdleChan
			w.status = Pending
			w.task = task
			r.RunningChan <- w
		case <-r.pauseChan:
			return
		}
	}
}

func (r *GuardPool) run() {
	go r.runStatusChan()
	go r.runWorkCore()
	go r.runTaskCore()
}

func (r *GuardPool) Start() {
	r.run()
}

func (r *GuardPool) Submit(task TaskFunc) {
	r.TaskFuncChan <- task
}

func (r *GuardPool) AsyncSubmit(task TaskFunc) {
	go func() {
		r.TaskFuncChan <- task
	}()
}

func (r *GuardPool) Release() {

}

func (r *GuardPool) Pause() {
	r.pauseChan <- struct{}{}
	r.running = false
}

func (r *GuardPool) Recover() {

	if !r.running {
		r.running = true
		go r.runTaskCore()
	}

}

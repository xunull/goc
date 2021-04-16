package routine_pool

import (
	"github.com/xunull/goc/easy"
	"sync"
)

type TaskFunc func()

const (
	Running = iota
	Pending
	Over
	Idle
)

type GoWorker struct {
	id       string
	taskFunc TaskFunc
	status   int
}

type RoutinePool struct {
	count        int
	RunningSheet map[string]*GoWorker
	IdleSheet    map[string]*GoWorker
	RunningChan  chan *GoWorker
	IdleChan     chan *GoWorker
	StatusChan   chan *GoWorker
	TaskFuncChan chan TaskFunc
	idMarker     *easy.IdMarker
	pauseChan    chan struct{}
	stopChan     chan struct{}
	running      bool
	runningCount int
	stopped      bool
	mu           sync.Mutex
}

func NewPool(count int) *RoutinePool {
	pool := &RoutinePool{
		RunningSheet: make(map[string]*GoWorker, count),
		IdleSheet:    make(map[string]*GoWorker, count),
		RunningChan:  make(chan *GoWorker, count),
		IdleChan:     make(chan *GoWorker, count),
		StatusChan:   make(chan *GoWorker, count),
		TaskFuncChan: make(chan TaskFunc, count),
		count:        count,
		idMarker:     &easy.IdMarker{},
		stopChan:     make(chan struct{}, 3),
	}

	pool.prepareWorker()
	return pool
}

func (r *RoutinePool) prepareWorker() {
	for i := 0; i < r.count; i++ {
		w := &GoWorker{
			id: r.idMarker.GetNewWorkerId(),
		}
		r.IdleChan <- w
	}
}

func (r *RoutinePool) GetRunningCount() int {
	return r.runningCount
}

func (r *RoutinePool) runStatusChan() {
	for {
		select {
		case w := <-r.StatusChan:
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
		case <-r.stopChan:
			return
		}
	}
}

func (r *RoutinePool) runWorkCore() {
	for {
		select {
		case w := <-r.RunningChan:
			if w.status == Pending {
				r.StatusChan <- w
			} else {
				go func(t *GoWorker) {
					t.taskFunc()
					t.taskFunc = nil
					t.status = Over
					r.StatusChan <- t
				}(w)
			}
		case <-r.stopChan:
			return
		}
	}

}

func (r *RoutinePool) runTaskChan() {
	for {
		select {
		case task := <-r.TaskFuncChan:
			w := <-r.IdleChan
			w.status = Pending
			w.taskFunc = task
			r.RunningChan <- w
		case <-r.pauseChan:
			return
		case <-r.stopChan:
			return
		}
	}
}

func (r *RoutinePool) run() {
	go r.runStatusChan()
	go r.runWorkCore()
	go r.runTaskChan()
}

func (r *RoutinePool) Start() {
	r.run()
}

func (r *RoutinePool) Submit(task TaskFunc) {
	go func() {
		r.TaskFuncChan <- task
	}()
}

func (r *RoutinePool) Release() {

}

func (r *RoutinePool) Pause() {
	r.pauseChan <- struct{}{}
	r.running = false
}

func (r *RoutinePool) Recover() {
	if !r.running {
		r.running = true
		go r.runTaskChan()
	}
}

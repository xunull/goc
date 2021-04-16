package routine_pool

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

type RoutinePool struct {
	count        int
	RunningSheet map[string]*GoWorker
	IdleSheet    map[string]*GoWorker
	RunningChan  chan *GoWorker
	IdleChan     chan *GoWorker
	StatusChan   chan *GoWorker
	TaskFuncChan chan TaskFunc
	idMarker     *IdMarker
	pauseChan    chan struct{}
	running      bool
	runningCount int
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
		idMarker:     &IdMarker{},
	}

	pool.prepareWorker()
	return pool
}

func (r *RoutinePool) prepareWorker() {
	for i := 0; i < r.count; i++ {
		w := &GoWorker{
			id: r.idMarker.getNewWorkerId(),
		}
		r.IdleChan <- w
	}
}

func (r *RoutinePool) runStatusChan() {
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

func (r *RoutinePool) GetRunningCount() int {
	return r.runningCount
}

func (r *RoutinePool) runWorkCore() {

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

func (r *RoutinePool) runTaskCore() {
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

func (r *RoutinePool) run() {
	go r.runStatusChan()
	go r.runWorkCore()
	go r.runTaskCore()
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
		go r.runTaskCore()
	}

}

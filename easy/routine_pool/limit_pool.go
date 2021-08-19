package routine_pool

import "github.com/xunull/goc/easy"

type LimitRoutinePool struct {
	limit      int
	doneCount  int
	noticeFunc func()
	*RoutinePool
}

func NewLimitPool(count int, limit int, notice func()) *LimitRoutinePool {
	if limit == 0 {
		panic("LimitRoutinePool Limit must not zero")
	}
	pool := RoutinePool{
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
	lrp := &LimitRoutinePool{
		limit:      limit,
		noticeFunc: notice,
	}
	lrp.RoutinePool = &pool
	return lrp
}

func (r *LimitRoutinePool) runStatusChan() {
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

				r.doneCount += 1

				go func() {
					r.noticeFunc()
				}()

				if r.doneCount >= r.limit {
					for i := 0; i < 3; i++ {
						r.stopChan <- struct{}{}
					}
				}
			}
		case <-r.stopChan:
			return
		}
	}
}

package traverse

import (
	"github.com/cheggaaa/pb/v3"
	"sync"
)

// the map[string]bool, "" item is self
// the name like path,name must unique
type Sheet map[string]bool

type WorkSheet struct {
	Sheet              Sheet
	Over               bool
	overChan           chan struct{}
	lock               *sync.RWMutex
	addOverFlag        bool
	overLock           *sync.RWMutex
	WithProgressBarOut bool
	pBar               *pb.ProgressBar
	Count              int64
	OverCount          int64
}

func (w *WorkSheet) StartProgressBar() {
	//w.WithProgressBarOut = true
	//w.pBar = pb.StartNew(0)
}

func NewWorkSheet() *WorkSheet {
	return &WorkSheet{
		Sheet:    make(map[string]bool),
		overChan: make(chan struct{}),
		lock:     &sync.RWMutex{},
		overLock: &sync.RWMutex{},
	}
}

func (w *WorkSheet) AddOver() {
	w.overLock.Lock()
	defer w.overLock.Unlock()
	w.lock.Lock()
	defer w.lock.Unlock()

	w.addOverFlag = true
	if w.Over {
		return
	}

	done := true
	for _, v := range w.Sheet {
		if !v {
			done = false
			break
		}
	}
	if done {
		w.Over = true
		w.overChan <- struct{}{}
	}
}

func (w *WorkSheet) IsOver() bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	done := true
	for _, v := range w.Sheet {
		if !v {
			done = false
			break
		}
	}
	return done
}

func (w *WorkSheet) Wait() {
	if w.Over {
		return
	}
	<-w.overChan
}

func (w *WorkSheet) Add(name string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.Count += 1

	//if w.WithProgressBarOut {
	//	w.pBar.SetTotal(w.Count)
	//}

	if _, ok := w.Sheet[name]; !ok {
		w.Sheet[name] = false
	}
}

func (w *WorkSheet) TargetOver(name string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.OverCount += 1

	//if w.WithProgressBarOut {
	//	w.pBar.SetCurrent(w.OverCount)
	//}

	w.Sheet[name] = true
	w.overLock.Lock()
	defer w.overLock.Unlock()
	if w.addOverFlag {

		if w.Over {
			return
		}

		done := true
		for _, v := range w.Sheet {
			if !v {
				done = false
				break
			}
		}
		if done {
			w.Over = true
			w.overChan <- struct{}{}

			//if w.WithProgressBarOut {
			//	w.pBar.Finish()
			//}

		}
	}

}

//func (w *WorkSheet) AddSub(name string, sub string) {
//	w.lock.Lock()
//	defer w.lock.Unlock()
//
//	if _, ok := w.Sheet[name]; ok {
//		if _, ok := w.Sheet[name][sub]; !ok {
//			w.Sheet[name][sub] = false
//		}
//	}
//	w.Add(sub)
//}
//
//func (w *WorkSheet) AddSubList(name string, subs []string) {
//	w.lock.Lock()
//	defer w.lock.Unlock()
//
//	if _, ok := w.Sheet[name]; ok {
//		for _, sub := range subs {
//			if _, ok := w.Sheet[name][sub]; !ok {
//				w.Sheet[name][sub] = false
//			}
//			w.Add(sub)
//		}
//	}
//}

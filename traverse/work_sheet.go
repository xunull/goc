package traverse

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
)

// the map[string]bool, "" item is self
// the name like path,name must unique
type Sheet map[string]bool

const (
	ItemAddStat = iota
	ItemOverStat
)

type itemStat struct {
	Name string
	Stat int
}

type WorkSheet struct {
	Sheet              Sheet
	Over               bool
	overChan           chan struct{}
	lock               *sync.RWMutex
	traverseOverFlag   bool
	overLock           *sync.RWMutex
	WithProgressBarOut bool
	pBar               *pb.ProgressBar
	Count              int64
	OverCount          int64
	statChan           chan itemStat
	cond               *sync.Cond // 使用条件变量替代 busy-wait
}

func NewWorkSheet() *WorkSheet {
	mutex := &sync.Mutex{}
	ws := &WorkSheet{
		Sheet:    make(map[string]bool),
		overChan: make(chan struct{}),
		lock:     &sync.RWMutex{},
		overLock: &sync.RWMutex{},
		statChan: make(chan itemStat, 2048),
		cond:     sync.NewCond(mutex),
	}
	go ws.runStatChanHandle()
	return ws
}

func (w *WorkSheet) ItemAdd(path string) {
	w.statChan <- itemStat{
		Name: path,
		Stat: ItemAddStat,
	}
}

func (w *WorkSheet) ItemOver(path string) {
	w.statChan <- itemStat{
		Name: path,
		Stat: ItemOverStat,
	}
}

func (w *WorkSheet) runStatChanHandle() {
	for stat := range w.statChan {
		w.cond.L.Lock()
		if stat.Stat == ItemAddStat {
			w.Count += 1

			if _, ok := w.Sheet[stat.Name]; !ok {
				w.Sheet[stat.Name] = false
			}

		} else {
			w.OverCount += 1
			w.Sheet[stat.Name] = true
		}
		// 通知等待的协程
		w.cond.Broadcast()
		w.cond.L.Unlock()
	}
}

func (w *WorkSheet) TraverseOver() {
	w.traverseOverFlag = true
	if w.Over {
		return
	}

	// 使用条件变量替代 busy-wait
	go func() {
		w.cond.L.Lock()
		for w.Count != w.OverCount {
			w.cond.Wait() // 等待信号，不消耗 CPU
		}
		w.Over = true
		w.cond.L.Unlock()
		w.overChan <- struct{}{}
	}()
}

func (w *WorkSheet) IsOver() bool {
	return w.Count == w.OverCount
}

func (w *WorkSheet) Wait() {
	if w.Over {
		return
	}
	<-w.overChan
}

// ---------------------------------------------------------------------------------------------------------------------

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

package trend_prediction

import "sync"

type DataSource struct {
	Source  chan interface{}
	KeyName string
}

type AlertFunc func(...interface{})

type AlertSubscribe struct {
	silence    bool
	inhibition bool
	AlertFunc  AlertFunc
}

func (s *AlertSubscribe) Silence() {
	s.silence = true
}

func (s *AlertSubscribe) Recover() {
	s.silence = false
}

func (s *AlertSubscribe) Inhibition() {
	s.inhibition = true
}

// ---------------------------------------------------------------------------------------------------------------------

type TrendPrediction struct {
	dataSource   *DataSource
	running      bool
	option       *option
	subscribers  []*AlertSubscribe
	subLocker    *sync.RWMutex
	lock         *sync.Mutex
	watchingData *dataInspect
}

func (t *TrendPrediction) getOption(opts ...Option) {
	if t.option == nil {
		t.option = getDefaultOption()
	}

	for _, f := range opts {
		f(t.option)
	}
}

func NewTrendPrediction(source *DataSource, opts ...Option) *TrendPrediction {
	tp := TrendPrediction{
		dataSource: source,
		subLocker:  &sync.RWMutex{},
		lock:       &sync.Mutex{},
	}
	tp.getOption(opts...)
	return &tp
}

func (t *TrendPrediction) Watching() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if !t.running {
		go t.startWatch()
	}
}

func (t *TrendPrediction) RegisterSubscribe(sub *AlertSubscribe) {
	t.subLocker.Lock()
	defer t.subLocker.Unlock()
	t.subscribers = append(t.subscribers, sub)
}

package trend_prediction

import (
	"container/list"
	"github.com/xunull/goc/enhance/slicex"
	"reflect"
	"time"
)

type trendData struct {
	value    int64
	ti       time.Time
	data     interface{}
	lastData *trendData
}

type dataInspect struct {
	leftTime  time.Time
	rightTime time.Time
	Length    int
	queue     *list.List
	dm        map[*trendData]*trendData
	dataChan  chan *trendData
	option    *option
	alertChan chan *trendData
}

func (di *dataInspect) checkTime() {
	now := time.Now()

	if !di.leftTime.IsZero() && now.Sub(di.leftTime) > time.Duration(di.option.WindowTimeEscape)*time.Second {
		di.PopFirst()
		di.checkTime()
	}
	return
}

func (di *dataInspect) PopFirst() {
	first := di.queue.Front()
	if first != nil {
		di.queue.Remove(first)
		temp := first.Value.(*trendData)
		delete(di.dm, temp)
		di.Length -= 1

		first = di.queue.Front()
		if first != nil {
			temp := first.Value.(*trendData)
			di.leftTime = temp.ti
		} else {
			di.leftTime = time.Time{}
		}

	} else {
		di.leftTime = time.Time{}
	}

}

func (di *dataInspect) inspect() {
	for data := range di.dataChan {

		if di.Length >= di.option.WindowLength {
			di.PopFirst()
		}

		di.checkTime()

		di.Length += 1
		di.queue.PushBack(data)
		di.dm[data] = data
		di.rightTime = data.ti

		di.calculate()
	}
}

func (di *dataInspect) calculate() {
	if di.option.CheckIncreasing {
		arr := make([]int64, 0, di.queue.Len())
		for e := di.queue.Front(); e != nil; e = e.Next() {
			temp := e.Value.(*trendData)
			arr = append(arr, temp.value)
		}
		if slicex.IsIncreasingInt64(arr) {
			b := di.queue.Back().Value.(*trendData)
			last := di.queue.Front().Value.(*trendData)
			b.lastData = last
			di.alertChan <- b
		}
	} else if di.option.CheckDescending {

	} else if di.option.CheckAverage {

	}
}

// ---------------------------------------------------------------------------------------------------------------------
func (t *TrendPrediction) initWatchingData() *dataInspect {
	wd := &dataInspect{
		dm:        make(map[*trendData]*trendData),
		dataChan:  make(chan *trendData, 512),
		queue:     list.New(),
		option:    t.option,
		alertChan: make(chan *trendData),
	}

	go t.runIncreasingAlert(wd.alertChan)
	go wd.inspect()

	return wd
}

func (t *TrendPrediction) startWatch() {
	t.lock.Lock()
	t.running = true
	t.lock.Unlock()

	t.watchingData = t.initWatchingData()

	for data := range t.dataSource.Source {
		t.processOne(data)
	}

	t.running = false
}

func (t *TrendPrediction) processOne(data interface{}) {

	ty := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
		v = v.Elem()
	}

	if _, found := ty.FieldByName(t.dataSource.KeyName); found {
		num := v.FieldByName(t.dataSource.KeyName).Int()

		t.watchingData.dataChan <- &trendData{
			value: num,
			ti:    time.Now(),
			data:  data,
		}
	}
}

func (t *TrendPrediction) runIncreasingAlert(ch <-chan *trendData) {
	for data := range ch {
		t.subLocker.Lock()
		for _, sub := range t.subscribers {
			sub.AlertFunc(data.data, data.lastData.data)
		}
		t.subLocker.Unlock()
	}
}

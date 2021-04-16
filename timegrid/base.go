package timegrid

import (
	"sync"
	"time"
)

type BoxTimeGrid struct {
	capacity        int
	tickerCallback  func()
	Store           map[int64][]interface{}
	BeatInterval    int
	mu              sync.Mutex
	timeTicker      *time.Ticker
	currentDataList []interface{}
}

func NewTimeGrid(capacity int) *BoxTimeGrid {
	tg := &BoxTimeGrid{
		capacity:   capacity,
		timeTicker: time.NewTicker(time.Second),
	}

	return tg
}

func (b *BoxTimeGrid) run() {
	go func() {
		select {
		case <-b.timeTicker.C:
		default:

		}
	}()
}

func (b *BoxTimeGrid) Put(data interface{}) {
	tu := time.Now().Unix()

	if dl, ok := b.Store[tu]; ok {
		dl = append(dl, data)
	} else {
		b.Store[tu] = make([]interface{}, 0, 1024)
		b.Store[tu] = append(b.Store[tu], data)
	}
}

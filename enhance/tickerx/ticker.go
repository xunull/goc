package tickerx

import (
	"github.com/xunull/goc/commonx"
	"time"
)

type SimpleFunc func()

func TickerSomeAndWaitSignal(interval int, funcs ...SimpleFunc) {
	for _, f := range funcs {
		go func(target SimpleFunc) {
			ticker := time.NewTicker(time.Second * time.Duration(interval))
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					target()
				}
			}
		}(f)
	}
	commonx.QuitWatch()
}

func TickerAndWaitSignal(interval int, f SimpleFunc) {

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(interval))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			}
		}
	}()

	commonx.QuitWatch()

}

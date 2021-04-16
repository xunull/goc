package easy

import (
	"fmt"
	"github.com/xunull/goc/enhance/timex"
	"sync"
)

type IdMarker struct {
	curIndex int64
	mu       sync.Mutex
}

func (m *IdMarker) GetNewWorkerId() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.curIndex += 1
	return fmt.Sprintf("GoWorker-%d-%s", m.curIndex, timex.GetYYYYMMDDHHMMSS())
}

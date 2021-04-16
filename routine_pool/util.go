package routine_pool

import (
	"fmt"
	"github.com/xunull/goc/enhance/timex"
)

type IdMarker struct {
	curIndex int64
}

func (m *IdMarker) getNewWorkerId() string {
	return fmt.Sprintf("GoWorker-%d-%s", m.curIndex, timex.GetYYYYMMDDHHMMSS())
}

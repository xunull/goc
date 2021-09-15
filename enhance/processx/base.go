package processx

import (
	"github.com/shirou/gopsutil/process"
	"os"
)

func GetParentName() (string, error) {
	id := os.Getppid()

	p, err := process.NewProcess(int32(id))
	if err != nil {
		return "", err
	}
	return p.Name()
}

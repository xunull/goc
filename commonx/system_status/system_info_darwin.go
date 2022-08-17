//go:build darwin

// Package system_status
// ------------------------------------------------------------------------
// Modified from gin-vue-admin (https://github.com/flipped-aurora/gin-vue-admin)
// ------------------------------------------------------------------------
package system_status

import (
	"runtime"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type SystemStatus struct {
	Os   Os   `json:"os"`
	Cpu  Cpu  `json:"cpu"`
	Ram  Ram  `json:"ram"`
	Disk Disk `json:"disk"`
}

type Os struct {
	GOOS         string `json:"goos"`
	NumCPU       int    `json:"numCpu"`
	Compiler     string `json:"compiler"`
	GoVersion    string `json:"goVersion"`
	NumGoroutine int    `json:"numGoroutine"`
}

type Cpu struct {
	Cpus  []float64 `json:"cpus"`
	Cores int       `json:"cores"`
}

type Ram struct {
	UsedMB      int `json:"usedMb"`
	TotalMB     int `json:"totalMb"`
	UsedPercent int `json:"usedPercent"`
}

type Disk struct {
	UsedMB      int `json:"usedMb"`
	UsedGB      int `json:"usedGb"`
	TotalMB     int `json:"totalMb"`
	TotalGB     int `json:"totalGb"`
	UsedPercent int `json:"usedPercent"`
}

func New() (s *SystemStatus, err error) {
	s = &SystemStatus{}
	s.Os = s.getOsStatus()
	if s.Cpu, err = s.getCpuStatus(); err != nil {
		return s, err
	}
	if s.Ram, err = s.getRamStatus(); err != nil {
		return s, err
	}
	if s.Disk, err = s.getDiskStatus(); err != nil {
		return s, err
	}
	return s, nil
}

func (s *SystemStatus) getOsStatus() (o Os) {
	o.GOOS = runtime.GOOS
	o.NumCPU = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.NumGoroutine = runtime.NumGoroutine()
	return o
}

func (s *SystemStatus) getCpuStatus() (c Cpu, err error) {
	return Cpu{}, nil
}

func (s *SystemStatus) getRamStatus() (r Ram, err error) {
	return Ram{}, nil
}

func (s *SystemStatus) getDiskStatus() (d Disk, err error) {
	return Disk{}, nil
}

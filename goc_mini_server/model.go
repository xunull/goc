package goc_mini_server

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/xunull/goc/commonx/system_status"
)

type ServerStatus struct {
	SystemStatus *system_status.SystemStatus
	Host         string
	Port         int
	RouteCount   int
	RouteTable   string
}

func (s ServerStatus) Output() {
	color.Yellow(fmt.Sprintf("Host: %s, Port: %d", s.Host, s.Port))
	color.Yellow(fmt.Sprintf("RouteCount: %d", s.RouteCount))
	color.Yellow(s.RouteTable)
}

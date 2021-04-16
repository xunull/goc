package commonx

import "github.com/xunull/goc/commonx/system_status"

func GetServerStatus() (*system_status.SystemStatus, error) {
	return system_status.New()
}

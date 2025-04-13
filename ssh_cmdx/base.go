package ssh_cmdx

import (
	"os"
	"strings"
)

type SSH struct {
	User     string `json:"user,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	PkName   string `json:"pkName,omitempty"`
	PkData   string `json:"pkData,omitempty"`
	Pk       string `json:"pk,omitempty"`
	PkPasswd string `json:"pkPasswd,omitempty"`
	Port     uint16 `json:"port,omitempty"`
}

func PathIsExist(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return os.IsExist(err)
	}
	return true
}

func GetHostIPAndPortOrDefault(host, Default string) (string, string) {
	if !strings.ContainsRune(host, ':') {
		return host, Default
	}
	split := strings.Split(host, ":")
	return split[0], split[1]
}

func GetSSHHostIPAndPort(host string) (string, string) {
	return GetHostIPAndPortOrDefault(host, "22")
}

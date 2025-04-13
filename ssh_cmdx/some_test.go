package ssh_cmdx

import (
	"context"
	"gopkg.in/yaml.v2"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const load1 = `

`

func loadYaml() map[string][]string {
	var load map[string][]string
	err := yaml.Unmarshal([]byte(load1), &load)
	if err != nil {
		panic(err)
	}
	return load
}

var testLoad = loadYaml()

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-1]
}

func getssh(param string) (Interface, error) {
	if strings.Contains(param, " ") {
		temp := strings.Split(param, " ")
		_, port, user, passwd := temp[0], temp[1], temp[2], temp[3]
		ss := &SSH{
			Pk: "/Users/quincy/.ssh/ansible",
		}
		if port != "-" {
			num, _ := strconv.ParseUint(port, 10, 16)
			ss.Port = uint16(num)
		}
		if user != "-" {
			ss.User = user
		}
		if passwd != "-" {
			ss.Passwd = passwd
		}
		return NewFromSSH(ss, false)

	} else {
		return NewFromSSH(&SSH{
			User: "root",
			Pk:   "/Users/quincy/.ssh/ansible",
		}, false)

	}
}

func getClientAndIp(param string) (Interface, string, error) {
	var client Interface
	var err error
	var ip string

	if strings.Contains(param, " ") {
		temp := strings.Split(param, " ")
		ip = temp[0]
		client, err = getssh(param)
	} else {
		ip = param
		client, err = getssh(param)
	}
	return client, ip, err
}

func TestCmdAsyncWithContext(t *testing.T) {
	funcName := getFuncName()
	t.Logf("funcName: %s", funcName)

	params := testLoad[funcName]
	for _, param := range params {
		t.Logf("param: %s", param)

		client, ip, err := getClientAndIp(param)

		if err != nil {
			t.Errorf("failed to create ssh client: %v", err)
		}

		err = client.CmdAsyncWithContext(context.Background(), ip, "hostname")
		if err != nil {
			t.Errorf("failed to CmdAsyncWithContext %s: %v", param, err)
		} else {
			t.Logf("[%s]  CmdAsyncWithContext success", param)
		}
	}
}

func TestCmdAsync(t *testing.T) {
	funcName := getFuncName()
	t.Logf("funcName: %s", funcName)

	params := testLoad[funcName]
	for _, param := range params {
		t.Logf("param: %s", param)

		client, ip, err := getClientAndIp(param)

		err = client.CmdAsync(ip, "hostname")
		if err != nil {
			t.Errorf("failed to CmdAsync %s: %v", param, err)
		} else {
			t.Logf("[%s]  CmdAsync success", param)
		}

	}
}

func TestCmd(t *testing.T) {
	funcName := getFuncName()
	t.Logf("funcName: %s", funcName)

	params := testLoad[funcName]
	for _, param := range params {
		t.Logf("param: %s", param)

		client, ip, err := getClientAndIp(param)

		if err != nil {
			t.Errorf("failed to create ssh client: %v", err)
		}

		out, err := client.Cmd(ip, "hostname")
		if err != nil {
			t.Errorf("failed to Cmd %s: %v", param, err)
		} else {
			t.Logf("[%s]  Cmd output: %s", param, string(out))
		}
	}
}

func TestCmdToString(t *testing.T) {
	funcName := getFuncName()
	t.Logf("funcName: %s", funcName)

	params := testLoad[funcName]
	for _, param := range params {
		t.Logf("param: %s", param)

		client, ip, err := getClientAndIp(param)

		if err != nil {
			t.Errorf("failed to create ssh client: %v", err)
		}

		out, err := client.CmdToString(ip, "hostname", "")
		if err != nil {
			t.Errorf("failed to CmdToString %s: %v", param, err)
		} else {
			t.Logf("[%s]  CmdToString output: %s", param, out)
		}
	}
}

func TestConnectToHost(t *testing.T) {
	funcName := getFuncName()
	t.Logf("funcName: %s", funcName)

	params := testLoad[funcName]
	for _, param := range params {
		t.Logf("param: %s", param)

		client, ip, err := getClientAndIp(param)
		if err != nil {
			t.Errorf("failed to create ssh client: %v", err)
		}
		err = client.Ping(ip)
		if err != nil {
			t.Errorf("failed to ping %s: %v", param, err)
		} else {
			t.Logf("ping %s success", param)
		}
	}
}

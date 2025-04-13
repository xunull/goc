package ssh_cmdx

import "context"

type Interface interface {
	// Copy copy local file to remote
	// scp -r /tmp root@192.168.0.2:/root/tmp => Copy("192.168.0.2","tmp","/root/tmp")
	// skip checksum if env DO_NOT_CHECKSUM=true
	Copy(host, src, dst string) error
	// Fetch fetch remote file to local
	// scp -r root@192.168.0.2:/remote/path/file /local/path/file => Fetch("192.168.0.2","/remote/path/file", "/local/path/file",)
	Fetch(host, src, dst string) error
	// CmdAsync exec commands on remote host asynchronously
	CmdAsync(host string, cmds ...string) error
	CmdAsyncWithContext(ctx context.Context, host string, cmds ...string) error
	// Cmd exec command on remote host, and return combined standard output and standard error
	Cmd(host, cmd string) ([]byte, error)
	// CmdToString exec command on remote host, and return spilt standard output by separator and standard error
	CmdToString(host, cmd, spilt string) (string, error)
	Ping(host string) error
}

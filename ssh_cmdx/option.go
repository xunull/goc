// Copyright © 2023 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh_cmdx

import (
	"github.com/xunull/goc/file_path"
	"net"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
)

type Option struct {
	stdout            bool
	sudo              bool
	user              string
	password          string
	privateKey        string
	rawPrivateKeyData string
	passphrase        string
	timeout           time.Duration
	hostKeyCallback   ssh.HostKeyCallback
}

const (
	defaultUsername = "root"
)

func NewOption() *Option {
	homedir, _ := os.UserHomeDir()
	getSSHFile := func(filenames ...string) string {
		for _, fn := range filenames {
			absPath := path.Join(homedir, ".ssh", fn)
			if ok, err := file_path.PathExists(absPath); ok && err == nil {
				return absPath
			}
		}
		return ""
	}
	opt := &Option{
		user:       defaultUsername,
		privateKey: getSSHFile("id_rsa", "id_dsa"),
		timeout:    10 * time.Second,
		hostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return opt
}

type OptionFunc func(*Option)

func WithSudoEnable(b bool) OptionFunc {
	return func(o *Option) {
		o.sudo = b
	}
}

func WithStdoutEnable(b bool) OptionFunc {
	return func(o *Option) {
		o.stdout = b
	}
}

func WithUsername(u string) OptionFunc {
	return func(o *Option) {
		o.user = u
	}
}

func WithPassword(p string) OptionFunc {
	return func(o *Option) {
		o.password = p
	}
}

func WithRawPrivateKeyDataAndPhrase(raw, passphrase string) OptionFunc {
	return func(o *Option) {
		o.rawPrivateKeyData = raw
		o.passphrase = passphrase
	}
}

func WithPrivateKeyAndPhrase(pk, passphrase string) OptionFunc {
	return func(o *Option) {
		o.privateKey = pk
		o.passphrase = passphrase
	}
}

func WithTimeout(timeout time.Duration) OptionFunc {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return func(o *Option) {
		o.timeout = timeout
	}
}

func WithHostKeyCallback(fn ssh.HostKeyCallback) OptionFunc {
	return func(o *Option) {
		o.hostKeyCallback = fn
	}
}

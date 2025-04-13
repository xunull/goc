// Copyright © 2021 Alibaba Group Holding Ltd.
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

import "golang.org/x/crypto/ssh"

func New(opt *Option, opts ...OptionFunc) (*Client, error) {
	return newFromOptions(opt, opts...)
}

func newFromOptions(opt *Option, opts ...OptionFunc) (*Client, error) {
	if opt == nil {
		opt = NewOption()
	}
	for i := range opts {
		opts[i](opt)
	}

	config := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: defaultCiphers,
		},
		User:            opt.user,
		Timeout:         opt.timeout,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: opt.hostKeyCallback,
	}
	if len(opt.password) > 0 {
		config.Auth = append(config.Auth, ssh.Password(opt.password))
	}
	if len(opt.rawPrivateKeyData) > 0 {
		signer, err := parsePrivateKey([]byte(opt.rawPrivateKeyData), []byte(opt.passphrase))
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	} else if len(opt.privateKey) > 0 {
		if !PathIsExist(opt.privateKey) {
			Debug("not trying to parse private key file cause it's not exists")
		} else {
			signer, err := parsePrivateKeyFile(opt.privateKey, opt.passphrase)
			if err != nil {
				return nil, err
			}
			config.Auth = append(config.Auth, ssh.PublicKeys(signer))
		}
	}
	return &Client{ClientConfig: config, Option: opt}, nil
}

func newOptionFromSSH(ssh *SSH, isStdout bool) *Option {
	opts := []OptionFunc{
		WithStdoutEnable(isStdout),
	}
	if len(ssh.User) > 0 {
		opts = append(opts, WithUsername(ssh.User))
	}
	if len(ssh.Passwd) > 0 {
		opts = append(opts, WithPassword(ssh.Passwd))
	}
	if len(ssh.Pk) > 0 {
		opts = append(opts, WithPrivateKeyAndPhrase(ssh.Pk, ssh.PkPasswd))
	}
	if len(ssh.PkData) > 0 {
		opts = append(opts, WithRawPrivateKeyDataAndPhrase(ssh.PkData, ssh.PkPasswd))
	}
	if ssh.User != "" && ssh.User != defaultUsername {
		opts = append(opts, WithSudoEnable(true))
	}

	opt := NewOption()
	for i := range opts {
		opts[i](opt)
	}
	return opt
}

func NewFromSSH(ssh *SSH, isStdout bool) (Interface, error) {
	return New(newOptionFromSSH(ssh, isStdout))
}

func MustNewClient(ssh *SSH, isStdout bool) Interface {
	client, err := NewFromSSH(ssh, isStdout)
	if err != nil {
		Fatal("failed to create ssh client: %v", err)
	}
	return client
}

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

import (
	"context"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

var (
	defaultMaxRetry         = 5
	defaultExecutionTimeout = 300 * time.Second
)

func RegisterFlags(fs *pflag.FlagSet) {
	fs.IntVar(&defaultMaxRetry, "max-retry", defaultMaxRetry, "define max num of ssh retry times")
	fs.DurationVar(&defaultExecutionTimeout, "execution-timeout", defaultExecutionTimeout, "timeout setting of command execution")
}

// GetTimeoutContext create a context.Context with default timeout
// default execution timeout in sealos is just fine, if you want to customize the timeout setting,
// you must invoke the `RegisterFlags` function above.
func GetTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultExecutionTimeout)
}

type Client struct {
	*ssh.ClientConfig
	*Option
}

var _ Interface = &Client{}

var defaultCiphers = []string{
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	"chacha20-poly1305@openssh.com",
	"aes128-gcm@openssh.com",
	"arcfour256", "arcfour128",
	"aes128-cbc", "aes192-cbc", "aes256-cbc",
	"3des-cbc",
}

// WaitReady wait for ssh ready
func WaitReady(client Interface, _ int, hosts ...string) error {
	eg, _ := errgroup.WithContext(context.Background())
	for i := range hosts {
		host := hosts[i]
		eg.Go(func() (err error) {
			return client.Ping(host)
		})
	}
	return eg.Wait()
}

// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package runtime

import (
	"fmt"
	"os"
	"os/exec"
)

// loadServiceProxyEnv sets the proxy variables of this process to those of the collector's systemd unit.
// Without systemd or the unit, such as in a container, the current environment is used.
func (c connectionTest) loadServiceProxyEnv() error {
	out, err := exec.Command("systemctl", "show", collectorServiceUnit,
		"--property=LoadState", "--property=Environment", "--property=EnvironmentFiles").Output()
	if err != nil {
		// systemctl is missing or systemd is not running
		c.printf("systemd is not available, using the current environment")
		return nil
	}

	info, err := parseSystemctlShow(string(out))
	if err != nil {
		return err
	}
	if info.loadState != "loaded" {
		c.printf("systemd unit %s not found, using the current environment", collectorServiceUnit)
		return nil
	}

	if os.Geteuid() != 0 {
		return fmt.Errorf("must be run as root to read the environment of systemd unit %s, re-run with sudo", collectorServiceUnit)
	}

	env, err := serviceProxyEnv(info, os.ReadFile)
	if err != nil {
		return err
	}
	c.printf("Using the proxy environment of systemd unit %s", collectorServiceUnit)
	return c.applyServiceProxyEnv(env)
}

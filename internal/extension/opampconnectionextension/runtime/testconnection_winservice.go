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

package runtime

import (
	"slices"
	"strings"
)

// windowsServiceProxyEnv returns the proxy variables a Windows service runs with. The service manager starts
// services with the system environment, and the service's Environment value overrides it. Windows environment
// variable names are case insensitive, so names are upper cased.
func windowsServiceProxyEnv(systemEnv map[string]string, serviceEnv []string) map[string]string {
	env := map[string]string{}
	set := func(name, value string) {
		name = strings.ToUpper(name)
		if slices.Contains(proxyEnvVars, name) || slices.Contains(noProxyEnvVars, name) {
			env[name] = value
		}
	}

	for name, value := range systemEnv {
		set(name, value)
	}
	for _, entry := range serviceEnv {
		if name, value, ok := strings.Cut(entry, "="); ok {
			set(name, value)
		}
	}
	return env
}

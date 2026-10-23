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

//go:build windows

package runtime

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	// collectorServiceName is the Windows service installed by the MSI
	collectorServiceName = "observiq-otel-collector"

	// collectorServiceKey is the registry key of the Windows service. Its Environment value holds the
	// variables set for the service.
	collectorServiceKey = `SYSTEM\CurrentControlSet\Services\` + collectorServiceName

	// systemEnvironmentKey is the registry key of the system environment variables
	systemEnvironmentKey = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
)

// loadServiceProxyEnv sets the proxy variables of this process to those of the collector's Windows service.
// Without the service, the current environment is used.
func (c connectionTest) loadServiceProxyEnv() error {
	serviceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, collectorServiceKey, registry.QUERY_VALUE)
	if errors.Is(err, registry.ErrNotExist) {
		c.printf("Windows service %s not found, using the current environment", collectorServiceName)
		return nil
	}
	if err != nil {
		return registryError(collectorServiceKey, err)
	}
	defer serviceKey.Close()

	serviceEnv, _, err := serviceKey.GetStringsValue("Environment")
	switch {
	case errors.Is(err, registry.ErrNotExist):
	case errors.Is(err, registry.ErrUnexpectedType):
		return fmt.Errorf(`the Environment value of HKLM\%s must be of type REG_MULTI_SZ`, collectorServiceKey)
	case err != nil:
		return registryError(collectorServiceKey, err)
	}

	systemEnv, err := systemProxyEnv()
	if err != nil {
		return err
	}

	c.printf("Using the proxy environment of Windows service %s", collectorServiceName)
	return c.applyServiceProxyEnv(windowsServiceProxyEnv(systemEnv, serviceEnv))
}

// systemProxyEnv returns the proxy variables set in the system environment
func systemProxyEnv() (map[string]string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, systemEnvironmentKey, registry.QUERY_VALUE)
	if err != nil {
		return nil, registryError(systemEnvironmentKey, err)
	}
	defer key.Close()

	env := map[string]string{}
	// Registry value names are case insensitive
	for _, name := range []string{"HTTPS_PROXY", "HTTP_PROXY", "NO_PROXY"} {
		value, valtype, err := key.GetStringValue(name)
		if errors.Is(err, registry.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, registryError(systemEnvironmentKey, err)
		}
		if valtype == registry.EXPAND_SZ {
			if value, err = registry.ExpandString(value); err != nil {
				return nil, fmt.Errorf("expand %s: %w", name, err)
			}
		}
		env[name] = value
	}
	return env, nil
}

func registryError(key string, err error) error {
	if errors.Is(err, windows.ERROR_ACCESS_DENIED) {
		return fmt.Errorf(`failed to read HKLM\%s: %w, re-run as Administrator`, key, err)
	}
	return fmt.Errorf(`failed to read HKLM\%s: %w`, key, err)
}

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
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// collectorServiceUnit is the systemd unit installed by the Linux packages
const collectorServiceUnit = "observiq-otel-collector.service"

// systemdUnitInfo is the subset of `systemctl show` output used to determine the service environment
type systemdUnitInfo struct {
	loadState        string
	environment      []string
	environmentFiles []systemdEnvironmentFile
}

// systemdEnvironmentFile is an EnvironmentFile= setting of a systemd unit
type systemdEnvironmentFile struct {
	path         string
	ignoreErrors bool
}

// parseSystemctlShow parses the output of `systemctl show --property=LoadState --property=Environment --property=EnvironmentFiles`.
// systemd has already merged the unit and its drop-ins.
func parseSystemctlShow(out string) (systemdUnitInfo, error) {
	var info systemdUnitInfo
	for _, line := range strings.Split(out, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		switch key {
		case "LoadState":
			info.loadState = value
		case "Environment":
			env, err := splitSystemdEnvironment(value)
			if err != nil {
				return info, err
			}
			info.environment = env
		case "EnvironmentFiles":
			// Each file is listed on its own line as "<path> (ignore_errors=<yes|no>)"
			path, opts := value, ""
			if i := strings.LastIndex(value, " (ignore_errors="); i >= 0 {
				path, opts = value[:i], value[i:]
			}
			info.environmentFiles = append(info.environmentFiles, systemdEnvironmentFile{
				path:         path,
				ignoreErrors: strings.Contains(opts, "ignore_errors=yes"),
			})
		}
	}
	return info, nil
}

// splitSystemdEnvironment splits the Environment property into KEY=VALUE entries.
// systemctl separates entries with spaces and shell quotes entries that contain special characters.
func splitSystemdEnvironment(s string) ([]string, error) {
	var entries []string
	var current strings.Builder
	var quote rune
	inEntry, escaped := false, false

	for _, r := range s {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\' && quote != '\'':
			escaped = true
			inEntry = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote = r
			inEntry = true
		case r == ' ' || r == '\t':
			if inEntry {
				entries = append(entries, current.String())
				current.Reset()
				inEntry = false
			}
		default:
			current.WriteRune(r)
			inEntry = true
		}
	}

	if quote != 0 || escaped {
		return nil, errors.New("failed to parse systemd Environment property: unterminated quote or escape")
	}
	if inEntry {
		entries = append(entries, current.String())
	}
	return entries, nil
}

// parseEnvironmentFile parses a systemd EnvironmentFile into KEY=VALUE entries.
// Blank lines and lines starting with # or ; are ignored, and quotes around values are removed.
// Line continuations are not supported.
func parseEnvironmentFile(data []byte) []string {
	var entries []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		entries = append(entries, key+"="+value)
	}
	return entries
}

// serviceProxyEnv returns the proxy variables the service runs with. As in systemd, EnvironmentFile= settings
// override Environment= settings, and later files override earlier ones.
func serviceProxyEnv(info systemdUnitInfo, readFile func(string) ([]byte, error)) (map[string]string, error) {
	env := map[string]string{}
	set := func(entries []string) {
		for _, entry := range entries {
			name, value, ok := strings.Cut(entry, "=")
			if ok && (slices.Contains(proxyEnvVars, name) || slices.Contains(noProxyEnvVars, name)) {
				env[name] = value
			}
		}
	}

	set(info.environment)
	for _, file := range info.environmentFiles {
		data, err := readFile(file.path)
		if err != nil {
			if file.ignoreErrors {
				continue
			}
			return nil, fmt.Errorf("failed to read environment file %s: %w", file.path, err)
		}
		set(parseEnvironmentFile(data))
	}
	return env, nil
}

// applyServiceProxyEnv replaces the proxy variables of this process with the service's
func (c connectionTest) applyServiceProxyEnv(env map[string]string) error {
	for _, name := range slices.Concat(proxyEnvVars, noProxyEnvVars) {
		if _, ok := env[name]; !ok && os.Getenv(name) != "" {
			c.printf("Ignoring %s from the current environment, it is not set for the collector service", name)
		}
		if err := os.Unsetenv(name); err != nil {
			return fmt.Errorf("unset %s: %w", name, err)
		}
	}

	for name, value := range env {
		if err := os.Setenv(name, value); err != nil {
			return fmt.Errorf("set %s: %w", name, err)
		}
	}
	return nil
}

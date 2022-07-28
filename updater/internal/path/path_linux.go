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

package path

import (
	"path/filepath"

	"go.uber.org/zap"
)

// LinuxInstallDir is the install directory of the collector on linux.
const LinuxInstallDir = "/opt/observiq-otel-collector"

// InstallDir returns the filepath to the install directory
func InstallDir(_ *zap.Logger) (string, error) {
	return LinuxInstallDir, nil
}

// LogFile returns the full path to the log file for the updater
func LogFile(installDir string) string {
	return filepath.Join(installDir, "log", "updater.log")
}

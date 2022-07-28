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

import "go.uber.org/zap"

// DarwinInstallDir is the path to the install directory on Darwin.
const DarwinInstallDir = "/opt/observiq-otel-collector"

// InstallDir returns the filepath to the install directory
func InstallDir(_ *zap.Logger) (string, error) {
	return DarwinInstallDir, nil
}

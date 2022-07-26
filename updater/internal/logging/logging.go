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

package logging

import (
	"fmt"
	"os"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"go.uber.org/zap"
)

// NewLogger returns a new logger, that logs to the log directory relative to installDir.
// It deletes the previous log file, as well.
func NewLogger(installDir string) (*zap.Logger, error) {
	logFile := path.LogFile(installDir)

	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{
		logFile,
	}

	err := os.RemoveAll(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to remove previous log file: %w", err)
	}

	prodLogger, err := conf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return prodLogger, nil
}

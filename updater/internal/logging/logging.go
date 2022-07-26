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

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(installDir string, level zapcore.Level) (*zap.Logger, error) {
	prodConf := zap.NewProductionConfig()

	prodLogger, err := prodConf.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core(installDir)
	}))

	if err != nil {
		return nil, fmt.Errorf("failed to make logger: %w", err)
	}

	return prodLogger, nil
}

func core(installDir string) zapcore.Core {
	logger := &lumberjack.Logger{
		Filename:   path.LogFile(installDir),
		MaxSize:    10,
		MaxBackups: 3,
	}

	return zapcore.NewCore(encoder(), zapcore.AddSync(logger), zapcore.DebugLevel)
}

func encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func LevelFromString(levelStr string) (zapcore.Level, error) {
	var l zapcore.Level = zapcore.DebugLevel
	err := l.Set(levelStr)
	return l, err
}

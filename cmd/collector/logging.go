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

package main

import (
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// logConfig represents the options used when configuring
// a zap logger.
type logConfig struct {
	path     string
	level    zapcore.Level
	maxSize  uint
	maxCount uint
	maxAge   uint
	compress bool
}

// opts returns zap log options.
func (l *logConfig) opts() []zap.Option {
	if l.path == "" {
		return nil
	}

	switch runtime.GOOS {
	case "windows":
		l.path = "winfile:///" + filepath.ToSlash(l.path)
	default:
		l.path = filepath.ToSlash(l.path)
	}

	logCore := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.MessageKey = "message"
		encoderConfig.CallerKey = "caller"
		encoderConfig.LevelKey = "level"
		encoderConfig.TimeKey = "timestamp"

		writer := &lumberjack.Logger{
			Filename:   l.path,
			MaxSize:    int(l.maxSize),
			MaxBackups: int(l.maxCount),
			MaxAge:     int(l.maxAge),
			Compress:   l.compress,
		}

		return zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			l.level,
		)
	})

	return []zap.Option{logCore}
}

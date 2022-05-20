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
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const (
	// FileOutput is an output option for logging to a file.
	FileOutput string = "file"

	// StdOutput is an output option for logging to stdout.
	StdOutput string = "stdout"
)

// LoggerConfig is the configuration of a logger.
type LoggerConfig struct {
	Output string             `yaml:"output"`
	Level  zapcore.Level      `yaml:"level"`
	File   *lumberjack.Logger `yaml:"file"`
}

// NewLoggerConfig returns a logger config. If configPath is not
// set, stdout logging will be enabled.
func NewLoggerConfig(configPath string) (*LoggerConfig, error) {
	conf := &LoggerConfig{
		Output: StdOutput,
		Level:  zapcore.InfoLevel,
	}

	if configPath == "" {
		return conf, nil
	}

	// If the file doesn't exist, we just return the default config.
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return conf, nil
	} else if err != nil {
		return nil, err
	}

	confBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(confBytes, &conf); err != nil {
		return nil, err
	}

	if conf.File != nil {
		// Expand optional environment variables in file path
		conf.File.Filename = os.ExpandEnv(conf.File.Filename)
	}

	return conf, nil
}

// Options returns the LoggerConfig's zap logging options.
func (l *LoggerConfig) Options() ([]zap.Option, error) {
	core, err := l.core()
	if err != nil {
		return nil, err
	}

	opt := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core
	})

	return []zap.Option{opt}, nil
}

// core returns the logging core specified in the config.
// An unknown output will return a nop core.
func (l *LoggerConfig) core() (zapcore.Core, error) {
	switch l.Output {
	case FileOutput:
		return zapcore.NewCore(newEncoder(), zapcore.AddSync(l.File), l.Level), nil
	case StdOutput:
		return zapcore.NewCore(newEncoder(), zapcore.Lock(os.Stdout), l.Level), nil
	default:
		return nil, fmt.Errorf("unrecognized output type: %s", l.Output)
	}
}

func newEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

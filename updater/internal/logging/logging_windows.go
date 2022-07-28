package logging

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"go.uber.org/zap"
)

var registerSinkOnce = &sync.Once{}

// NewLogger returns a new logger, that logs to the log directory relative to installDir.
// It deletes the previous log file, as well.
// NewLogger must only be called once, at the start of the program.
func NewLogger(installDir string) (*zap.Logger, error) {
	// On windows, absolute paths do not work for zap's default sink, so we must register it.
	// see: https://github.com/uber-go/zap/issues/621
	var err error
	registerSinkOnce.Do(func() {
		err = zap.RegisterSink("winfile", newWinFileSink)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to registed windows file sink: %w", err)
	}

	logFile := path.LogFile(installDir)

	err = os.RemoveAll(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to remove previous log file: %w", err)
	}

	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{
		"winfile:///" + logFile,
	}

	prodLogger, err := conf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return prodLogger, nil
}

// Windows requires a special sink, so that we may properly parse the file path
// See: https://github.com/uber-go/zap/issues/621
func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
}

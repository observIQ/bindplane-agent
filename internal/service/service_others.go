//go:build !windows

package service

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// RunService runs the given service, calling its start and stop functions.
func RunService(logger *zap.Logger, rSvc RunnableService) error {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stopSignal)

	return runServiceInteractive(logger, stopSignal, rSvc)
}

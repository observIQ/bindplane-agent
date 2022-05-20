//go:build windows

package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
)

// The following constants specify error codes for the service.
// See https://docs.microsoft.com/en-us/windows/win32/debug/system-error-codes--1000-1299-
const (
	statusCodeInvalidServiceCommand = uint32(1052)
	statusCodeServiceException      = uint32(1064)
	statusCodeInvalidServiceName    = uint32(1213)
)

func RunService(logger *zap.Logger, rSvc RunnableService) error {
	isService, err := checkIsService()
	if err != nil {
		return fmt.Errorf("failed checking if running as service: %w", err)
	}

	if isService {
		// Service name doesn't need to be specified when directly run by the service manager.
		return svc.Run("", newWindowsServiceHandler(logger, rSvc))
	} else {
		stopSignal := make(chan os.Signal, 1)
		signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stopSignal)

		return runServiceInteractive(logger, stopSignal, rSvc)
	}
}

// windowsServiceHandler implements svc.Handler
type windowsServiceHandler struct {
	svc    RunnableService
	logger *zap.Logger
}

// newWindowsServiceHandler creates a new windowsServiceHandler, which implements svc.Handler
func newWindowsServiceHandler(logger *zap.Logger, svc RunnableService) *windowsServiceHandler {
	return &windowsServiceHandler{
		svc:    svc,
		logger: logger,
	}
}

// Execute handles the Windows service event loop.
func (sh *windowsServiceHandler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	if len(args) == 0 {
		// Service name is the first argument, and must be provided to open the event log for service logs.
		return false, statusCodeInvalidServiceName
	}

	s <- svc.Status{State: svc.StartPending}

	startupTimeoutCtx, startupCancel := context.WithTimeout(context.Background(), startTimeout)
	defer startupCancel()

	err := sh.svc.Start(startupTimeoutCtx)
	if err != nil {
		sh.logger.Error("Failed to start service", zap.Error(err))
		return false, statusCodeServiceException
	}

	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	for {
		select {
		case req := <-r:
			switch req.Cmd {
			case svc.Interrogate:
				s <- req.CurrentStatus
			case svc.Stop, svc.Shutdown:
				err := sh.shutdown(s)
				if err != nil {
					sh.logger.Error("Failed during service shutdown", zap.Error(err))
					return false, statusCodeServiceException
				}

				return false, 0
			default:
				sh.logger.Error("Got unexpected service command", zap.Uint32("command", uint32(req.Cmd)))
				err := sh.shutdown(s)
				if err != nil {
					sh.logger.Error("Failed during service shutdown", zap.Error(err))
					return false, statusCodeServiceException
				}

				return false, statusCodeInvalidServiceCommand
			}
		case err := <-sh.svc.Error():
			sh.logger.Error("Got unexpected service error", zap.Error(err))

			sh.shutdown(s)

			if err != nil {
				sh.logger.Error("Failed during service shutdown", zap.Error(err))
			}

			return false, statusCodeServiceException
		}
	}
}

func (sh windowsServiceHandler) shutdown(s chan<- svc.Status) error {
	s <- svc.Status{State: svc.StopPending}

	stopTimeoutCtx, stopCancel := context.WithTimeout(context.Background(), stopTimeout)
	defer stopCancel()

	err := sh.svc.Stop(stopTimeoutCtx)

	s <- svc.Status{State: svc.Stopped}

	return err
}

// checkIsService returns whether the current process is running as a Windows service.
func checkIsService() (bool, error) {
	// NO_WINDOWS_SERVICE may be set non-zero to override the service detection logic.
	if value, present := os.LookupEnv("NO_WINDOWS_SERVICE"); present && value != "0" {
		return true, nil
	}

	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		return false, fmt.Errorf("failed to determine if we are running in an windows service: %w", err)
	}

	return isWindowsService, nil
}

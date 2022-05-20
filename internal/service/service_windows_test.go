package service

import (
	"errors"
	"testing"
	"time"

	"github.com/observiq/observiq-otel-collector/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
)

func TestWindowsServiceHandler(t *testing.T) {
	t.Run("Normal start/stop", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(make(chan error)))
		rSvc.On("Stop", mock.Anything).Return(nil)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		changeChan <- svc.ChangeRequest{
			Cmd:           svc.Interrogate,
			CurrentStatus: svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown},
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for interrogate response")
		}

		changeChan <- svc.ChangeRequest{
			Cmd: svc.Stop,
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, uint32(0), statusCode, "status code was not 0")
	})

	t.Run("Start fails", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		startError := errors.New("Failed to start service")

		rSvc.On("Start", mock.Anything).Return(startError)
		rSvc.On("Error").Return((<-chan error)(make(chan error)))
		rSvc.On("Stop", mock.Anything).Return(nil)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, statusCodeServiceException, statusCode, "status code was not ServiceException")
	})

	t.Run("Unexpected service error", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		svcErr := errors.New("service unexpectedly failed")
		errChan := make(chan error, 1)
		errChan <- svcErr

		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(errChan))
		rSvc.On("Stop", mock.Anything).Return(nil)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, statusCodeServiceException, statusCode, "status code was not ServiceException")
	})

	t.Run("Shutdown error", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}
		stopError := errors.New("Failed to start service")

		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(make(chan error)))
		rSvc.On("Stop", mock.Anything).Return(stopError)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		changeChan <- svc.ChangeRequest{
			Cmd:           svc.Interrogate,
			CurrentStatus: svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown},
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for interrogate response")
		}

		changeChan <- svc.ChangeRequest{
			Cmd: svc.Stop,
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, statusCodeServiceException, statusCode, "status code was not ServiceException")
	})

	t.Run("Shutdown error after unexpected error", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}
		stopError := errors.New("Failed to start service")
		svcErr := errors.New("service unexpectedly failed")
		errChan := make(chan error, 1)
		errChan <- svcErr

		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(errChan))
		rSvc.On("Stop", mock.Anything).Return(stopError)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, statusCodeServiceException, statusCode, "status code was not ServiceException")
	})

	t.Run("Unhandled command", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(make(chan error)))
		rSvc.On("Stop", mock.Anything).Return(nil)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		changeChan <- svc.ChangeRequest{
			Cmd:           svc.Interrogate,
			CurrentStatus: svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown},
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for interrogate response")
		}

		changeChan <- svc.ChangeRequest{
			Cmd: svc.DeviceEvent,
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, uint32(statusCodeInvalidServiceCommand), statusCode, "status code was not InvalidServiceCommand")
	})

	t.Run("Unhandled command with shutdown error", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		stopError := errors.New("Failed to start service")
		rSvc.On("Start", mock.Anything).Return(nil)
		rSvc.On("Error").Return((<-chan error)(make(chan error)))
		rSvc.On("Stop", mock.Anything).Return(stopError)

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{"service-name"}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StartPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to start pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service status change to running")
		}

		changeChan <- svc.ChangeRequest{
			Cmd:           svc.Interrogate,
			CurrentStatus: svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown},
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for interrogate response")
		}

		changeChan <- svc.ChangeRequest{
			Cmd: svc.DeviceEvent,
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.StopPending}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stop pending")
		}

		select {
		case status := <-statusChan:
			require.Equal(t, svc.Status{State: svc.Stopped}, status)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for status change to stopped")
		}

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, uint32(statusCodeServiceException), statusCode, "status code was not ServiceException")
	})

	t.Run("No service name", func(t *testing.T) {
		rSvc := &mocks.RunnableService{}

		svcHandler := newWindowsServiceHandler(zap.NewNop(), rSvc)

		changeChan := make(chan svc.ChangeRequest)
		statusChan := make(chan svc.Status, 6)
		svcHandlerDone := make(chan struct{})

		var isSvcSpecificStatus bool
		var statusCode uint32
		go func() {
			isSvcSpecificStatus, statusCode = svcHandler.Execute([]string{}, changeChan, statusChan)
			close(svcHandlerDone)
		}()

		select {
		case <-svcHandlerDone: // OK
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for service handler to return")
		}

		require.Equal(t, false, isSvcSpecificStatus, "status code marked as service specific")
		require.Equal(t, uint32(statusCodeInvalidServiceName), statusCode, "status code was not InvalidServiceName")
	})
}

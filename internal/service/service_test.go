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

package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStandaloneCollectorService(t *testing.T) {
	t.Run("Collector starts and stops normally", func(t *testing.T) {
		col := &mocks.Collector{}
		col.On("Run", context.Background()).Return(nil)
		col.On("Status").Return((<-chan *collector.Status)(make(chan *collector.Status)))
		col.On("Stop", mock.Anything).Return(nil)

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			err = srv.Start(context.Background())
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(time.Second):
			t.Fatalf("Start timed out")
		}

		require.NoError(t, err)
		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")

		stoppedChan := make(chan struct{})
		go func() {
			err = srv.Stop(context.Background())
			close(stoppedChan)
		}()

		select {
		case <-stoppedChan: // OK
		case <-time.After(time.Second):
			t.Fatalf("Stop timed out")
		}

		require.NoError(t, err)
	})

	t.Run("Collector.Run errors", func(t *testing.T) {
		col := &mocks.Collector{}
		runError := errors.New("run failed")

		col.On("Run", context.Background()).Return(runError)
		col.On("Status").Return((<-chan *collector.Status)(make(chan *collector.Status)))
		col.On("Stop", mock.Anything).Return(nil)

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			err = srv.Start(context.Background())
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(time.Second):
			t.Fatalf("Start timed out")
		}

		defer srv.Stop(context.Background())

		require.NoError(t, err)
		select {
		case err := <-srv.Error():
			require.Equal(t, runError, err)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for error")
		}
		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")
	})

	t.Run("Start context times out", func(t *testing.T) {
		col := &mocks.Collector{}

		col.On("Run", context.Background()).Run(func(args mock.Arguments) { time.Sleep(100 * time.Second) })
		col.On("Status").Return((<-chan *collector.Status)(make(chan *collector.Status)))
		col.On("Stop", mock.Anything).Return(nil)

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			err = srv.Start(ctx)
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(2 * time.Second):
			t.Fatalf("Start timed out")
		}

		defer srv.Stop(context.Background())

		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")
	})

	t.Run("Stop context times out", func(t *testing.T) {
		col := &mocks.Collector{}
		col.On("Run", context.Background()).Return(nil)
		col.On("Status").Return((<-chan *collector.Status)(make(chan *collector.Status)))
		col.On("Stop", mock.Anything).Run(func(args mock.Arguments) { time.Sleep(100 * time.Second) })

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			err = srv.Start(context.Background())
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(time.Second):
			t.Fatalf("Start timed out")
		}

		require.NoError(t, err)
		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")

		stoppedChan := make(chan struct{})
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			err = srv.Stop(ctx)
			close(stoppedChan)
		}()

		select {
		case <-stoppedChan: // OK
		case <-time.After(time.Second):
			t.Fatalf("Stop timed out")
		}

		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("Collector status has an error", func(t *testing.T) {
		col := &mocks.Collector{}
		colStatusErr := errors.New("Collector errored")
		colStatus := make(chan *collector.Status, 1)
		colStatus <- &collector.Status{
			Running: false,
			Err:     colStatusErr,
		}

		col.On("Run", context.Background()).Return(nil)
		col.On("Status").Return((<-chan *collector.Status)(colStatus))
		col.On("Stop", mock.Anything).Return(nil)

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			err = srv.Start(context.Background())
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(2 * time.Second):
			t.Fatalf("Start timed out")
		}

		require.NoError(t, err)

		defer srv.Stop(context.Background())

		select {
		case err := <-srv.Error():
			require.Equal(t, colStatusErr, err)
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for error")
		}

		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")
	})

	t.Run("Collector status is not running", func(t *testing.T) {
		col := &mocks.Collector{}
		colStatus := make(chan *collector.Status, 1)
		colStatus <- &collector.Status{
			Running: false,
		}

		col.On("Run", context.Background()).Return(nil)
		col.On("Status").Return((<-chan *collector.Status)(colStatus))
		col.On("Stop", mock.Anything).Return(nil)

		srv := NewStandaloneCollectorService(col)

		var err error
		startedChan := make(chan struct{})
		go func() {
			err = srv.Start(context.Background())
			close(startedChan)
		}()

		select {
		case <-startedChan: // OK
		case <-time.After(2 * time.Second):
			t.Fatalf("Start timed out")
		}

		require.NoError(t, err)

		defer srv.Stop(context.Background())

		select {
		case err := <-srv.Error():
			require.Contains(t, err.Error(), "collector unexpectedly stopped running")
		case <-time.After(time.Second):
			t.Fatalf("Timed out waiting for error")
		}

		require.Equal(t, 0, len(srv.Error()), "error channel has elements in it!")
	})
}

func TestRunServiceInteractive(t *testing.T) {
	t.Run("Normal start/stop", func(t *testing.T) {
		svc := &mocks.RunnableService{}

		svc.On("Start", mock.Anything).Return(nil)
		svc.On("Error").Return((<-chan error)(make(chan error)))
		svc.On("Stop", mock.Anything).Return(nil)

		stopSignal := make(chan os.Signal, 1)

		var err error
		svcDone := make(chan struct{})
		go func() {
			err = runServiceInteractive(zap.NewNop(), stopSignal, svc)
			close(svcDone)
		}()

		stopSignal <- os.Interrupt

		select {
		case <-svcDone: // OK
		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for service done")
		}

		require.NoError(t, err)
	})

	t.Run("Start fails", func(t *testing.T) {
		svc := &mocks.RunnableService{}

		startErr := errors.New("failed to start")

		svc.On("Start", mock.Anything).Return(startErr)

		stopSignal := make(chan os.Signal, 1)

		var err error
		svcDone := make(chan struct{})
		go func() {
			err = runServiceInteractive(zap.NewNop(), stopSignal, svc)
			close(svcDone)
		}()

		select {
		case <-svcDone: // OK
		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for service done")
		}

		require.Error(t, err)
		require.ErrorIs(t, err, startErr)
	})

	t.Run("Service errors", func(t *testing.T) {
		svc := &mocks.RunnableService{}

		svcErr := errors.New("service unexpectedly failed")
		errChan := make(chan error, 1)
		errChan <- svcErr

		svc.On("Start", mock.Anything).Return(nil)
		svc.On("Error").Return((<-chan error)(errChan))
		svc.On("Stop", mock.Anything).Return(nil)

		stopSignal := make(chan os.Signal, 1)

		var err error
		svcDone := make(chan struct{})
		go func() {
			err = runServiceInteractive(zap.NewNop(), stopSignal, svc)
			close(svcDone)
		}()

		select {
		case <-svcDone: // OK
		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for service done")
		}

		require.Error(t, err)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("Stop errors", func(t *testing.T) {
		svc := &mocks.RunnableService{}

		stopErr := errors.New("Stop failed")

		svc.On("Start", mock.Anything).Return(nil)
		svc.On("Error").Return((<-chan error)(make(chan error)))
		svc.On("Stop", mock.Anything).Return(stopErr)

		stopSignal := make(chan os.Signal, 1)

		var err error
		svcDone := make(chan struct{})
		go func() {
			err = runServiceInteractive(zap.NewNop(), stopSignal, svc)
			close(svcDone)
		}()

		stopSignal <- os.Interrupt

		select {
		case <-svcDone: // OK
		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for service done")
		}

		require.Error(t, err)
		require.ErrorIs(t, err, stopErr)
	})

	t.Run("Stop errors after error returned", func(t *testing.T) {
		svc := &mocks.RunnableService{}

		stopErr := errors.New("Stop failed")
		svcErr := errors.New("service unexpectedly failed")
		errChan := make(chan error, 1)
		errChan <- svcErr

		svc.On("Start", mock.Anything).Return(nil)
		svc.On("Error").Return((<-chan error)(errChan))
		svc.On("Stop", mock.Anything).Return(stopErr)

		stopSignal := make(chan os.Signal, 1)

		var err error
		svcDone := make(chan struct{})
		go func() {
			err = runServiceInteractive(zap.NewNop(), stopSignal, svc)
			close(svcDone)
		}()

		select {
		case <-svcDone: // OK
		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for service done")
		}

		require.Error(t, err)
		require.ErrorIs(t, err, stopErr)
	})
}

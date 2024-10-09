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

package state

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/observiq/bindplane-agent/packagestate/mocks"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCollectorMonitorSetState(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Package not in current status",
			testFunc: func(*testing.T) {
				mockStateManger := mocks.NewMockStateManager(t)

				collectorMonitor := &CollectorMonitor{
					stateManager: mockStateManger,
					currentStatus: &protobufs.PackageStatuses{
						Packages: make(map[string]*protobufs.PackageStatus),
					},
				}

				err := collectorMonitor.SetState("my_package", protobufs.PackageStatusEnum_PackageStatusEnum_Installed, nil)
				assert.Error(t, err)
			},
		},
		{
			desc: "Sets Status no error",
			testFunc: func(*testing.T) {
				pgkName := "my_package"
				expectedStatus := &protobufs.PackageStatuses{
					Packages: map[string]*protobufs.PackageStatus{
						pgkName: {
							Name:                 pgkName,
							AgentHasVersion:      "1.0",
							AgentHasHash:         []byte("hash1"),
							ServerOfferedVersion: "1.2",
							ServerOfferedHash:    []byte("hash2"),
							Status:               protobufs.PackageStatusEnum_PackageStatusEnum_Installed,
						},
					},
				}

				mockStateManger := mocks.NewMockStateManager(t)
				mockStateManger.On("SaveStatuses", expectedStatus).Return(nil)

				collectorMonitor := &CollectorMonitor{
					stateManager: mockStateManger,
					currentStatus: &protobufs.PackageStatuses{
						Packages: map[string]*protobufs.PackageStatus{
							pgkName: {
								Name:                 pgkName,
								AgentHasVersion:      "1.0",
								AgentHasHash:         []byte("hash1"),
								ServerOfferedVersion: "1.2",
								ServerOfferedHash:    []byte("hash2"),
								Status:               protobufs.PackageStatusEnum_PackageStatusEnum_InstallPending,
							},
						},
					},
				}

				err := collectorMonitor.SetState("my_package", protobufs.PackageStatusEnum_PackageStatusEnum_Installed, nil)
				assert.NoError(t, err)
				assert.Equal(t, expectedStatus, collectorMonitor.currentStatus)
			},
		},
		{
			desc: "Sets Status w/error",
			testFunc: func(*testing.T) {
				pgkName := "my_package"
				statusErr := errors.New("some error")

				expectedStatus := &protobufs.PackageStatuses{
					Packages: map[string]*protobufs.PackageStatus{
						pgkName: {
							Name:                 pgkName,
							AgentHasVersion:      "1.0",
							AgentHasHash:         []byte("hash1"),
							ServerOfferedVersion: "1.2",
							ServerOfferedHash:    []byte("hash2"),
							Status:               protobufs.PackageStatusEnum_PackageStatusEnum_InstallFailed,
							ErrorMessage:         statusErr.Error(),
						},
					},
				}

				mockStateManger := mocks.NewMockStateManager(t)
				mockStateManger.On("SaveStatuses", expectedStatus).Return(nil)

				collectorMonitor := &CollectorMonitor{
					stateManager: mockStateManger,
					currentStatus: &protobufs.PackageStatuses{
						Packages: map[string]*protobufs.PackageStatus{
							pgkName: {
								Name:                 pgkName,
								AgentHasVersion:      "1.0",
								AgentHasHash:         []byte("hash1"),
								ServerOfferedVersion: "1.2",
								ServerOfferedHash:    []byte("hash2"),
								Status:               protobufs.PackageStatusEnum_PackageStatusEnum_InstallPending,
							},
						},
					},
				}

				err := collectorMonitor.SetState("my_package", protobufs.PackageStatusEnum_PackageStatusEnum_InstallFailed, statusErr)
				assert.NoError(t, err)
				assert.Equal(t, expectedStatus, collectorMonitor.currentStatus)
			},
		},
		{
			desc: "StateManager fails to save",
			testFunc: func(*testing.T) {
				pgkName := "my_package"
				expectedErr := errors.New("bad")
				expectedStatus := &protobufs.PackageStatuses{
					Packages: map[string]*protobufs.PackageStatus{
						pgkName: {
							Name:                 pgkName,
							AgentHasVersion:      "1.0",
							AgentHasHash:         []byte("hash1"),
							ServerOfferedVersion: "1.2",
							ServerOfferedHash:    []byte("hash2"),
							Status:               protobufs.PackageStatusEnum_PackageStatusEnum_Installed,
						},
					},
				}

				mockStateManger := mocks.NewMockStateManager(t)
				mockStateManger.On("SaveStatuses", expectedStatus).Return(expectedErr)

				collectorMonitor := &CollectorMonitor{
					stateManager: mockStateManger,
					currentStatus: &protobufs.PackageStatuses{
						Packages: map[string]*protobufs.PackageStatus{
							pgkName: {
								Name:                 pgkName,
								AgentHasVersion:      "1.0",
								AgentHasHash:         []byte("hash1"),
								ServerOfferedVersion: "1.2",
								ServerOfferedHash:    []byte("hash2"),
								Status:               protobufs.PackageStatusEnum_PackageStatusEnum_InstallPending,
							},
						},
					},
				}

				err := collectorMonitor.SetState("my_package", protobufs.PackageStatusEnum_PackageStatusEnum_Installed, nil)
				assert.ErrorIs(t, err, expectedErr)
				assert.Equal(t, expectedStatus, collectorMonitor.currentStatus)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestCollectorMonitorMonitorForSuccess(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Context is canceled",
			testFunc: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
				defer testServer.Close()
				port := testServer.Listener.Addr().(*net.TCPAddr).Port

				collectorMonitor := &CollectorMonitor{
					logger:          zap.NewNop(),
					healthCheckPort: port,
				}

				err := collectorMonitor.MonitorForSuccess(ctx)
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
		{
			desc: "Successful startup",
			testFunc: func(t *testing.T) {
				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
				defer testServer.Close()
				port := testServer.Listener.Addr().(*net.TCPAddr).Port

				collectorMonitor := &CollectorMonitor{
					logger:          zap.NewNop(),
					healthCheckPort: port,
				}

				err := collectorMonitor.MonitorForSuccess(context.Background())
				assert.NoError(t, err)
			},
		},
		{
			desc: "Unreachable agent",
			testFunc: func(t *testing.T) {
				collectorMonitor := &CollectorMonitor{
					logger:          zap.NewNop(),
					healthCheckPort: 12345,
				}

				err := collectorMonitor.MonitorForSuccess(context.Background())
				assert.ErrorContains(t, err, "failed to reach agent after 3 attempts")
			},
		},
		{
			desc: "Agent returns bad status",
			testFunc: func(t *testing.T) {
				testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusBadRequest) }))
				defer testServer.Close()
				port := testServer.Listener.Addr().(*net.TCPAddr).Port

				collectorMonitor := &CollectorMonitor{
					logger:          zap.NewNop(),
					healthCheckPort: port,
				}

				err := collectorMonitor.MonitorForSuccess(context.Background())
				assert.ErrorContains(t, err, fmt.Sprintf("health check on %d returned %d", port, http.StatusBadRequest))
			},
		},
		{
			desc: "Agent initially fails but then succeeds",
			testFunc: func(t *testing.T) {
				testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
				port := testServer.Listener.Addr().(*net.TCPAddr).Port
				err := testServer.Listener.Close()
				require.NoError(t, err)

				collectorMonitor := &CollectorMonitor{
					logger:          zap.NewNop(),
					healthCheckPort: port,
				}

				go func() {
					// simulate agent coming up late
					time.Sleep(time.Second * 5)
					l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
					require.NoError(t, err)
					testServer.Listener = l
					testServer.Start()
				}()

				err = collectorMonitor.MonitorForSuccess(context.Background())
				assert.NoError(t, err)
				testServer.Close()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

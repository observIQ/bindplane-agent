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

package report

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/observiq/observiq-otel-collector/internal/report/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManagerSetClient(t *testing.T) {
	testCases := []struct {
		desc        string
		client      Client
		expectedErr error
	}{
		{
			desc:        "Nil client",
			client:      nil,
			expectedErr: errors.New("client must not be nil"),
		},
		{
			desc:        "Successful set",
			client:      http.DefaultClient,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := &Manager{}
			err := m.SetClient(tc.client)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.client, m.client)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestManagerResetConfig(t *testing.T) {
	testCases := []struct {
		desc          string
		configPath    string
		mockSetupFunc func(*testing.T, *Manager)
		expectedErr   error
	}{
		{
			desc:          "bad config",
			configPath:    "./testdata/bad_config.yaml",
			mockSetupFunc: func(*testing.T, *Manager) {},
			expectedErr:   errors.New("failed to unmarshal config"),
		},
		{
			desc:          "Unknown reporter kind",
			configPath:    "./testdata/unknown_reporter.yaml",
			mockSetupFunc: func(*testing.T, *Manager) {},
			expectedErr:   errors.New("unrecognized reporter kind"),
		},
		{
			desc:          "Bad Snapshot",
			configPath:    "./testdata/bad_snapshot.yaml",
			mockSetupFunc: func(*testing.T, *Manager) {},
			expectedErr:   errors.New("failed to unmarshal Snapshot config"),
		},
		{
			desc:       "Valid config, reporter fails to stop",
			configPath: "./testdata/valid.yaml",
			mockSetupFunc: func(t *testing.T, m *Manager) {
				mockSnapshotReporter := mocks.NewMockReporter(t)
				mockSnapshotReporter.On("Stop", mock.Anything).Return(errors.New("bad"))
				mockSnapshotReporter.On("Kind").Return(snapShotKind)

				m.reporters[snapShotKind] = mockSnapshotReporter

			},
			expectedErr: errors.New("failed to stop"),
		},
		{
			desc:       "Valid config, reporter fails to Report",
			configPath: "./testdata/valid.yaml",
			mockSetupFunc: func(t *testing.T, m *Manager) {
				mockSnapshotReporter := mocks.NewMockReporter(t)
				mockSnapshotReporter.On("Stop", mock.Anything).Return(nil)
				mockSnapshotReporter.On("Report", mock.Anything).Return(errors.New("bad"))
				mockSnapshotReporter.On("Kind").Return(snapShotKind)

				m.reporters[snapShotKind] = mockSnapshotReporter

			},
			expectedErr: errors.New("failed to report"),
		},
		{
			desc:       "Valid config, no errors",
			configPath: "./testdata/valid.yaml",
			mockSetupFunc: func(t *testing.T, m *Manager) {
				mockSnapshotReporter := mocks.NewMockReporter(t)
				mockSnapshotReporter.On("Stop", mock.Anything).Return(nil)
				mockSnapshotReporter.On("Report", mock.Anything).Return(nil)

				m.reporters[snapShotKind] = mockSnapshotReporter

			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := &Manager{
				client:    http.DefaultClient,
				reporters: make(map[string]Reporter),
			}

			tc.mockSetupFunc(t, m)

			configData, err := os.ReadFile(tc.configPath)
			assert.NoError(t, err)

			err = m.ResetConfig(configData)
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManagerShutdown(t *testing.T) {
	testCases := []struct {
		desc          string
		mockSetupFunc func(*testing.T, *Manager)
		expectedErr   error
	}{
		{
			desc: "Reporter errors",
			mockSetupFunc: func(t *testing.T, m *Manager) {
				mockSnapshotReporter := mocks.NewMockReporter(t)
				mockSnapshotReporter.On("Stop", mock.Anything).Return(errors.New("bad"))
				mockSnapshotReporter.On("Kind").Return(snapShotKind)

				m.reporters[snapShotKind] = mockSnapshotReporter
			},
			expectedErr: errors.New("bad"),
		},
		{
			desc: "No errors",
			mockSetupFunc: func(t *testing.T, m *Manager) {
				mockSnapshotReporter := mocks.NewMockReporter(t)
				mockSnapshotReporter.On("Stop", mock.Anything).Return(nil)

				m.reporters[snapShotKind] = mockSnapshotReporter
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := &Manager{
				client:    http.DefaultClient,
				reporters: make(map[string]Reporter),
			}

			tc.mockSetupFunc(t, m)

			err := m.Shutdown(context.Background())
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

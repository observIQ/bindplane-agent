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

//go:build !windows

package observiq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewOthersUpdaterManager(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "New LinuxUpdaterManager",
			testFunc: func(t *testing.T) {
				tmpPath := "/tmp"
				logger := zap.NewNop()

				expected := &othersUpdaterManager{
					tmpPath: tmpPath,
					logger:  logger.Named("updater manager"),
				}

				actual := newUpdaterManager(logger, tmpPath)
				require.Equal(t, expected, actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

// We don't have a good way to unit test the happy path,
// which involves the entire collector being killed in the middle of this function
func TestStartAndMonitorUpdater(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Updater does not exist at path",
			testFunc: func(t *testing.T) {
				tmpDir := t.TempDir()
				updateManager := newUpdaterManager(zap.NewNop(), tmpDir)
				err := updateManager.StartAndMonitorUpdater()

				assert.ErrorContains(t, err, "no such file or directory")
			},
		},
		{
			desc: "Updater is not executable",
			testFunc: func(t *testing.T) {
				updaterName = "badupdater"
				updateManager := newUpdaterManager(zap.NewNop(), "./testdata")
				err := updateManager.StartAndMonitorUpdater()

				assert.ErrorContains(t, err, "updater had an issue while starting:")
			},
		},
		{
			desc: "Updater exits quickly",
			testFunc: func(t *testing.T) {
				updaterName = "quickupdater"
				updateManager := newUpdaterManager(zap.NewNop(), "./testdata")
				err := updateManager.StartAndMonitorUpdater()

				assert.EqualError(t, err, "updater failed to update collector")
			},
		},
		{
			desc: "Updater times out",
			testFunc: func(t *testing.T) {
				updaterName = "slowupdater"
				updateManager := newUpdaterManager(zap.NewNop(), "./testdata")
				err := updateManager.StartAndMonitorUpdater()

				assert.ErrorContains(t, err, "updater failed to update collector")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

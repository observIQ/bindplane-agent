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

package observiq

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/packagestate"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewPackagesStateProvider(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "New PackagesStateProvider",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				actual := newPackagesStateProvider(logger, "test.json")

				packagesStateProvider, ok := actual.(*packagesStateProvider)
				require.True(t, ok)

				assert.Equal(t, logger, packagesStateProvider.logger)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestAllPackagesHash(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				actual, err := p.AllPackagesHash()

				assert.Nil(t, actual)
				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestSetAllPackagesHash(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				err := p.SetAllPackagesHash([]byte("hash"))

				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestPackages(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				actual, err := p.Packages()

				assert.Nil(t, actual)
				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestPackageState(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				actual, err := p.PackageState("name")

				assert.Equal(t, types.PackageState{}, actual)
				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestSetPackageState(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				err := p.SetPackageState("name", types.PackageState{})

				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestCreatePackage(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				err := p.CreatePackage("name", protobufs.PackageAvailable_TopLevelPackage)

				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestFileContentHash(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				actual, err := p.FileContentHash("name")

				assert.Nil(t, actual)
				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestUpdateContent(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}
				var r io.Reader

				err := p.UpdateContent(context.TODO(), "name", r, []byte("hash"))

				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestDeletePackage(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Not Implemented",
			testFunc: func(t *testing.T) {
				logger := zap.NewNop()
				p := &packagesStateProvider{
					logger: logger,
				}

				err := p.DeletePackage("name")

				assert.ErrorContains(t, err, "method not implemented")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestLastReportedStatuses(t *testing.T) {
	pkgName := packagestate.CollectorPackageName
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "No Existing file",
			testFunc: func(t *testing.T) {
				noExistJSON := "garbage.json"
				logger := zap.NewNop()
				p := newPackagesStateProvider(logger, noExistJSON)

				actual, err := p.LastReportedStatuses()

				assert.NoError(t, err)
				assert.Nil(t, actual.ServerProvidedAllPackagesHash)
				assert.Equal(t, "", actual.ErrorMessage)
				assert.Equal(t, 1, len(actual.Packages))
				assert.Equal(t, pkgName, actual.Packages[pkgName].GetName())
				assert.Equal(t, version.Version(), actual.Packages[pkgName].GetAgentHasVersion())
				assert.Nil(t, actual.Packages[pkgName].GetAgentHasHash())
				assert.Equal(t, "", actual.Packages[pkgName].GetServerOfferedVersion())
				assert.Nil(t, actual.Packages[pkgName].GetServerOfferedHash())
				assert.Equal(t, protobufs.PackageStatus_Installed, actual.Packages[pkgName].GetStatus())
				assert.Equal(t, "", actual.Packages[pkgName].GetErrorMessage())
			},
		},
		{
			desc: "Load Error",
			testFunc: func(t *testing.T) {
				badJSON := "testdata/package_statuses_bad.json"
				logger := zap.NewNop()
				p := newPackagesStateProvider(logger, badJSON)

				actual, err := p.LastReportedStatuses()

				assert.Error(t, err)
				assert.Nil(t, actual.ServerProvidedAllPackagesHash)
				assert.Equal(t, "", actual.ErrorMessage)
				assert.Equal(t, 1, len(actual.Packages))
				assert.Equal(t, pkgName, actual.Packages[pkgName].GetName())
				assert.Equal(t, version.Version(), actual.Packages[pkgName].GetAgentHasVersion())
				assert.Nil(t, actual.Packages[pkgName].GetAgentHasHash())
				assert.Equal(t, "", actual.Packages[pkgName].GetServerOfferedVersion())
				assert.Nil(t, actual.Packages[pkgName].GetServerOfferedHash())
				assert.Equal(t, protobufs.PackageStatus_Installed, actual.Packages[pkgName].GetStatus())
				assert.Equal(t, "", actual.Packages[pkgName].GetErrorMessage())
			},
		},
		{
			desc: "Successful file read",
			testFunc: func(t *testing.T) {
				goodJSON := "testdata/package_statuses_good.json"
				pkgName := "package"
				agentVersion := "1.0"
				agentHash := []byte("hash1")
				serverVersion := "2.0"
				serverHash := []byte("hash2")
				errMsg := "bad"
				allHash := []byte("hash")
				allErrMsg := "whoops"
				logger := zap.NewNop()
				p := newPackagesStateProvider(logger, goodJSON)

				actual, err := p.LastReportedStatuses()

				assert.NoError(t, err)
				assert.Equal(t, allHash, actual.ServerProvidedAllPackagesHash)
				assert.Equal(t, allErrMsg, actual.ErrorMessage)
				assert.Equal(t, 1, len(actual.Packages))
				assert.Equal(t, pkgName, actual.Packages[pkgName].GetName())
				assert.Equal(t, agentVersion, actual.Packages[pkgName].GetAgentHasVersion())
				assert.Equal(t, agentHash, actual.Packages[pkgName].GetAgentHasHash())
				assert.Equal(t, serverVersion, actual.Packages[pkgName].GetServerOfferedVersion())
				assert.Equal(t, serverHash, actual.Packages[pkgName].GetServerOfferedHash())
				assert.Equal(t, protobufs.PackageStatus_InstallPending, actual.Packages[pkgName].GetStatus())
				assert.Equal(t, errMsg, actual.Packages[pkgName].GetErrorMessage())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestSetLastReportedStatuses(t *testing.T) {
	pkgName := "package"
	agentVersion := "1.0"
	agentHash := []byte("hash1")
	serverVersion := "2.0"
	serverHash := []byte("hash2")
	errMsg := "bad"
	allHash := []byte("hash")
	allErrMsg := "whoops"
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "New file",
			testFunc: func(t *testing.T) {
				tmpDir := t.TempDir()
				testJSON := filepath.Join(tmpDir, "test.json")
				logger := zap.NewNop()
				p := newPackagesStateProvider(logger, testJSON)

				provider, ok := p.(*packagesStateProvider)
				assert.True(t, ok)

				packages := map[string]*protobufs.PackageStatus{
					pkgName: {
						Name:                 pkgName,
						AgentHasVersion:      agentVersion,
						AgentHasHash:         agentHash,
						ServerOfferedVersion: serverVersion,
						ServerOfferedHash:    serverHash,
						Status:               protobufs.PackageStatus_InstallPending,
						ErrorMessage:         errMsg,
					},
				}
				packageStatuses := &protobufs.PackageStatuses{
					ServerProvidedAllPackagesHash: allHash,
					ErrorMessage:                  allErrMsg,
					Packages:                      packages,
				}

				err := p.SetLastReportedStatuses(packageStatuses)
				assert.NoError(t, err)

				actual, err := provider.packageStateManager.LoadStatuses()

				assert.NoError(t, err)
				assert.NoError(t, err)
				assert.Equal(t, allHash, actual.GetServerProvidedAllPackagesHash())
				assert.Equal(t, allErrMsg, actual.GetErrorMessage())
				assert.Equal(t, 1, len(actual.GetPackages()))
				assert.Equal(t, pkgName, actual.GetPackages()[pkgName].GetName())
				assert.Equal(t, agentVersion, actual.GetPackages()[pkgName].GetAgentHasVersion())
				assert.Equal(t, agentHash, actual.GetPackages()[pkgName].GetAgentHasHash())
				assert.Equal(t, serverVersion, actual.GetPackages()[pkgName].GetServerOfferedVersion())
				assert.Equal(t, serverHash, actual.GetPackages()[pkgName].GetServerOfferedHash())
				assert.Equal(t, protobufs.PackageStatus_InstallPending, actual.GetPackages()[pkgName].GetStatus())
				assert.Equal(t, errMsg, actual.GetPackages()[pkgName].GetErrorMessage())
			},
		},
		{
			desc: "Existing file",
			testFunc: func(t *testing.T) {
				tmpDir := t.TempDir()
				testJSON := filepath.Join(tmpDir, "test.json")
				os.WriteFile(testJSON, nil, 0600)

				logger := zap.NewNop()
				p := newPackagesStateProvider(logger, testJSON)
				provider, ok := p.(*packagesStateProvider)
				assert.True(t, ok)

				packages := map[string]*protobufs.PackageStatus{
					pkgName: {
						Name:                 pkgName,
						AgentHasVersion:      agentVersion,
						AgentHasHash:         agentHash,
						ServerOfferedVersion: serverVersion,
						ServerOfferedHash:    serverHash,
						Status:               protobufs.PackageStatus_InstallPending,
						ErrorMessage:         errMsg,
					},
				}
				packageStatuses := &protobufs.PackageStatuses{
					ServerProvidedAllPackagesHash: allHash,
					ErrorMessage:                  allErrMsg,
					Packages:                      packages,
				}

				err := provider.SetLastReportedStatuses(packageStatuses)
				assert.NoError(t, err)

				actual, err := provider.packageStateManager.LoadStatuses()
				assert.NoError(t, err)

				assert.NoError(t, err)
				assert.Equal(t, allHash, actual.GetServerProvidedAllPackagesHash())
				assert.Equal(t, allErrMsg, actual.GetErrorMessage())
				assert.Equal(t, 1, len(actual.GetPackages()))
				assert.Equal(t, pkgName, actual.GetPackages()[pkgName].GetName())
				assert.Equal(t, agentVersion, actual.GetPackages()[pkgName].GetAgentHasVersion())
				assert.Equal(t, agentHash, actual.GetPackages()[pkgName].GetAgentHasHash())
				assert.Equal(t, serverVersion, actual.GetPackages()[pkgName].GetServerOfferedVersion())
				assert.Equal(t, serverHash, actual.GetPackages()[pkgName].GetServerOfferedHash())
				assert.Equal(t, protobufs.PackageStatus_InstallPending, actual.GetPackages()[pkgName].GetStatus())
				assert.Equal(t, errMsg, actual.GetPackages()[pkgName].GetErrorMessage())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

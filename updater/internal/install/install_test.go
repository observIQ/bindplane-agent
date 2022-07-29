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

package install

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-otel-collector/updater/internal/action"
	rb_mocks "github.com/observiq/observiq-otel-collector/updater/internal/rollback/mocks"
	"github.com/observiq/observiq-otel-collector/updater/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestInstallArtifacts(t *testing.T) {
	t.Run("Installs artifacts correctly", func(t *testing.T) {
		outDir := t.TempDir()
		svc := mocks.NewService(t)
		rb := rb_mocks.NewActionAppender(t)

		installer := &Installer{
			latestDir:  filepath.Join("testdata", "example-install"),
			installDir: outDir,
			svc:        svc,
			logger:     zaptest.NewLogger(t),
		}

		outDirConfig := filepath.Join(outDir, "config.yaml")
		outDirLogging := filepath.Join(outDir, "logging.yaml")
		outDirManager := filepath.Join(outDir, "manager.yaml")

		err := os.WriteFile(outDirConfig, []byte("# The original config file"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(outDirLogging, []byte("# The original logging file"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(outDirManager, []byte("# The original manager file"), 0600)
		require.NoError(t, err)

		svc.On("Stop").Once().Return(nil)
		svc.On("Update").Once().Return(nil)
		svc.On("Start").Once().Return(nil)

		actions := []action.RollbackableAction{}
		rb.On("AppendAction", mock.Anything).Run(func(args mock.Arguments) {
			action := args.Get(0).(action.RollbackableAction)
			actions = append(actions, action)
		})

		err = installer.Install(rb)
		require.NoError(t, err)

		contentsEqual(t, outDirConfig, "# The original config file")
		contentsEqual(t, outDirManager, "# The original manager file")
		contentsEqual(t, outDirLogging, "# The original logging file")

		require.FileExists(t, filepath.Join(outDir, "test.txt"))
		require.DirExists(t, filepath.Join(outDir, "test-folder"))
		require.FileExists(t, filepath.Join(outDir, "test-folder", "another-test.txt"))

		contentsEqual(t, filepath.Join(outDir, "test.txt"), "This is a test file\n")
		contentsEqual(t, filepath.Join(outDir, "test-folder", "another-test.txt"), "This is a nested text file\n")

		copyTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test.txt"),
			filepath.Join(installer.installDir, "test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyTestTxtAction.FileCreated = true

		copyNestedTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test-folder", "another-test.txt"),
			filepath.Join(installer.installDir, "test-folder", "another-test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyNestedTestTxtAction.FileCreated = true

		require.Equal(t, []action.RollbackableAction{
			action.NewServiceStopAction(svc),
			copyNestedTestTxtAction,
			copyTestTxtAction,
			action.NewServiceUpdateAction(installer.logger, installer.tmpDir),
			action.NewServiceStartAction(svc),
		}, actions)
	})

	t.Run("Stop fails", func(t *testing.T) {
		outDir := t.TempDir()
		svc := mocks.NewService(t)
		rb := rb_mocks.NewActionAppender(t)
		installer := &Installer{
			latestDir:  filepath.Join("testdata", "example-install"),
			installDir: outDir,
			svc:        svc,
			logger:     zaptest.NewLogger(t),
		}

		svc.On("Stop").Once().Return(errors.New("stop failed"))

		err := installer.Install(rb)
		require.ErrorContains(t, err, "failed to stop service")
	})

	t.Run("Update fails", func(t *testing.T) {
		outDir := t.TempDir()
		svc := mocks.NewService(t)
		rb := rb_mocks.NewActionAppender(t)
		installer := &Installer{
			latestDir:  filepath.Join("testdata", "example-install"),
			installDir: outDir,
			svc:        svc,
			logger:     zaptest.NewLogger(t),
		}

		svc.On("Stop").Once().Return(nil)
		svc.On("Update").Once().Return(errors.New("uninstall failed"))

		actions := []action.RollbackableAction{}
		rb.On("AppendAction", mock.Anything).Run(func(args mock.Arguments) {
			action := args.Get(0).(action.RollbackableAction)
			actions = append(actions, action)
		})

		err := installer.Install(rb)
		require.ErrorContains(t, err, "failed to update service")
		copyTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test.txt"),
			filepath.Join(installer.installDir, "test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyTestTxtAction.FileCreated = true

		copyNestedTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test-folder", "another-test.txt"),
			filepath.Join(installer.installDir, "test-folder", "another-test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyNestedTestTxtAction.FileCreated = true

		require.Equal(t, []action.RollbackableAction{
			action.NewServiceStopAction(svc),
			copyNestedTestTxtAction,
			copyTestTxtAction,
		}, actions)
	})

	t.Run("Start fails", func(t *testing.T) {
		outDir := t.TempDir()
		svc := mocks.NewService(t)
		rb := rb_mocks.NewActionAppender(t)
		installer := &Installer{
			latestDir:  filepath.Join("testdata", "example-install"),
			installDir: outDir,
			svc:        svc,
			logger:     zaptest.NewLogger(t),
		}

		svc.On("Stop").Once().Return(nil)
		svc.On("Update").Once().Return(nil)
		svc.On("Start").Once().Return(errors.New("start failed"))

		actions := []action.RollbackableAction{}
		rb.On("AppendAction", mock.Anything).Run(func(args mock.Arguments) {
			action := args.Get(0).(action.RollbackableAction)
			actions = append(actions, action)
		})

		err := installer.Install(rb)
		require.ErrorContains(t, err, "failed to start service")

		copyTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test.txt"),
			filepath.Join(installer.installDir, "test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyTestTxtAction.FileCreated = true

		copyNestedTestTxtAction, err := action.NewCopyFileAction(
			installer.logger,
			filepath.Join("test-folder", "another-test.txt"),
			filepath.Join(installer.installDir, "test-folder", "another-test.txt"),
			installer.tmpDir,
		)
		require.NoError(t, err)
		copyNestedTestTxtAction.FileCreated = true

		require.Equal(t, []action.RollbackableAction{
			action.NewServiceStopAction(svc),
			copyNestedTestTxtAction,
			copyTestTxtAction,
			action.NewServiceUpdateAction(installer.logger, installer.tmpDir),
		}, actions)
	})

	t.Run("Latest dir does not exist", func(t *testing.T) {
		outDir := t.TempDir()
		svc := mocks.NewService(t)
		rb := rb_mocks.NewActionAppender(t)
		installer := &Installer{
			latestDir:  filepath.Join("testdata", "non-existent-dir"),
			installDir: outDir,
			svc:        svc,
			logger:     zaptest.NewLogger(t),
		}

		svc.On("Stop").Once().Return(nil)

		actions := []action.RollbackableAction{}
		rb.On("AppendAction", mock.Anything).Run(func(args mock.Arguments) {
			action := args.Get(0).(action.RollbackableAction)
			actions = append(actions, action)
		})

		err := installer.Install(rb)
		require.ErrorContains(t, err, "failed to install new files")

		require.Equal(t, []action.RollbackableAction{
			action.NewServiceStopAction(svc),
		}, actions)
	})
}

func contentsEqual(t *testing.T, path, expectedContents string) {
	t.Helper()

	contents, err := os.ReadFile(path)
	require.NoError(t, err)

	// Replace \r\n with \n to normalize for windows tests.
	contents = bytes.ReplaceAll(contents, []byte("\r\n"), []byte("\n"))
	require.Equal(t, []byte(expectedContents), contents)
}

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

// Package install handles installation of new collector artifacts
package install

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/observiq/bindplane-agent/updater/internal/action"
	"github.com/observiq/bindplane-agent/updater/internal/file"
	"github.com/observiq/bindplane-agent/updater/internal/path"
	"github.com/observiq/bindplane-agent/updater/internal/rollback"
	"github.com/observiq/bindplane-agent/updater/internal/service"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Installer is an interface that performs an Install of a new collector.
//
//go:generate mockery --name Installer --filename mock_installer.go --structname MockInstaller
type Installer interface {
	// Install installs new artifacts over the old ones.
	Install(rollback.Rollbacker, int) error
}

// archiveInstaller allows you to install files from latestDir into installDir,
// as well as update the service configuration using the "Install" method.
type archiveInstaller struct {
	latestDir  string
	installDir string
	backupDir  string
	svc        service.Service
	logger     *zap.Logger
}

// NewInstaller returns a new instance of an Installer.
func NewInstaller(logger *zap.Logger, installDir string, service service.Service) Installer {
	return &archiveInstaller{
		latestDir:  path.LatestDir(installDir),
		svc:        service,
		installDir: installDir,
		backupDir:  path.BackupDir(installDir),
		logger:     logger.Named("installer"),
	}
}

// Install installs the unpacked artifacts in latestDir to installDir,
// as well as installing the new service file using the installer's Service interface.
// It then starts the service.
func (i archiveInstaller) Install(rb rollback.Rollbacker, hcePort int) error {
	// If JMX jar exists outside of install directory, make sure that gets backed up
	if err := i.attemptSpecialJMXJarInstall(rb); err != nil {
		return fmt.Errorf("failed to process special JMX jar: %w", err)
	}

	// install files that go to installDirPath to their correct location,
	// excluding any config files (logging.yaml, config.yaml, manager.yaml)
	if err := installFiles(i.logger, i.latestDir, i.installDir, i.backupDir, rb); err != nil {
		return fmt.Errorf("failed to install new files: %w", err)
	}
	i.logger.Debug("Install artifacts copied")

	// translate manager.yaml into supervisor config
	if err := initializeSupervisorConfig(i.logger, i.installDir, i.backupDir, rb, hcePort); err != nil {
		return fmt.Errorf("failed to translate manager config into supervisor config: %w", err)
	}
	i.logger.Debug("Translated config files")

	// Update old service config to new service config
	if err := i.svc.Update(); err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	rb.AppendAction(action.NewServiceUpdateAction(i.logger, i.installDir))
	i.logger.Debug("Updated service configuration")

	// Start service
	if err := i.svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	rb.AppendAction(action.NewServiceStartAction(i.svc))
	i.logger.Debug("Service started")

	return nil
}

// installFiles moves the file tree rooted at inputPath to installDir,
// skipping configuration files. Appends CopyFileAction-s to the Rollbacker as it copies file.
func installFiles(logger *zap.Logger, inputPath, installDir, backupDir string, rb rollback.Rollbacker) error {
	err := filepath.WalkDir(inputPath, func(inPath string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			// if there was an error walking the directory, we want to bail out.
			return err
		case d.IsDir():
			// Skip directories, we'll create them when we get a file in the directory.
			return nil
		case skipConfigFiles(inPath):
			// Found a config file that we should skip copying.
			return nil
		}

		// We want the path relative to the directory we are walking in order to calculate where the file should be
		// mirrored in the destination directory.
		relPath, err := filepath.Rel(inputPath, inPath)
		if err != nil {
			return err
		}

		// use the relative path to get the outPath (where we should write the file), and
		// to get the out directory (which we will create if it does not exist).
		outPath := filepath.Join(installDir, relPath)
		outDir := filepath.Dir(outPath)

		if err := os.MkdirAll(outDir, 0750); err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}

		// We create the action record here, because we want to record whether the file exists or not before
		// we open the file (which will end up creating the file).
		cfa, err := action.NewCopyFileAction(logger, relPath, outPath, backupDir)
		if err != nil {
			return fmt.Errorf("failed to create copy file action: %w", err)
		}

		// Record that we are performing copying the file.
		// We record before we actually do the action here because the file may be partially written,
		// and we will want to roll that back if that is the case.
		rb.AppendAction(cfa)

		if err := file.CopyFileOverwrite(logger.Named("copy-file"), inPath, outPath); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk latest dir: %w", err)
	}

	return nil
}

func (i archiveInstaller) attemptSpecialJMXJarInstall(rb rollback.Rollbacker) error {
	jarPath := path.SpecialJMXJarFile(i.installDir)
	jarDirPath := path.SpecialJarDir(i.installDir)
	latestJarPath := path.LatestJMXJarFile(i.latestDir)
	_, err := os.Stat(jarPath)
	switch {
	case err == nil:
		if err := installFile(i.logger, latestJarPath, jarDirPath, i.backupDir, rb); err != nil {
			return fmt.Errorf("failed to install JMX jar from latest directory: %w", err)
		}
		// Just log this error as the worst case is that there will be two jars copied over
		if err = os.Remove(latestJarPath); err != nil {
			i.logger.Warn("Failed to remove JMX jar from latest directory", zap.Error(err))
		}
	case !errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("failed determine where currently installed JMX jar is: %w", err)
	}

	return nil
}

// installFile moves new file to output path.
// Appends CopyFileAction-s to the Rollbacker as it copies file.
func installFile(logger *zap.Logger, inPath, installDirPath, backupDirPath string, rb rollback.Rollbacker) error {
	baseInPath := filepath.Base(inPath)

	// use the relative path to get the outPath (where we should write the file), and
	// to get the out directory (which we will create if it does not exist).
	outPath := filepath.Join(installDirPath, baseInPath)
	outDir := filepath.Dir(outPath)

	if err := os.MkdirAll(outDir, 0750); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	// We create the action record here, because we want to record whether the file exists or not before
	// we open the file (which will end up creating the file).
	cfa, err := action.NewCopyFileAction(logger, baseInPath, outPath, backupDirPath)
	if err != nil {
		return fmt.Errorf("failed to create copy file action: %w", err)
	}

	// Record that we are performing copying the file.
	// We record before we actually do the action here because the file may be partially written,
	// and we will want to roll that back if that is the case.
	rb.AppendAction(cfa)

	if err := file.CopyFileOverwrite(logger.Named("copy-file"), inPath, outPath); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// skipConfigFiles returns true if the given path is a special config file.
// These files should not be overwritten.
func skipConfigFiles(path string) bool {
	var configFiles = []string{
		"config.yaml",
		"logging.yaml",
		"manager.yaml",
	}

	fileName := filepath.Base(path)

	for _, f := range configFiles {
		if fileName == f {
			return true
		}
	}

	return false
}

type Headers struct {
	Authorization string `yaml:"Authorization"`
	AgentIDFormat string `yaml:"X-Bindplane-Agent-Id-Format"`
}

type TLS struct {
	Insecure           bool `yaml:"insecure"`
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
}

type Server struct {
	Endpoint string  `yaml:"endpoint"`
	Headers  Headers `yaml:"headers"`
	TLS      TLS     `yaml:"tls"`
}

type Capabilities struct {
	AcceptsRemoteConfig bool `yaml:"accepts_remote_config"`
	ReportsRemoteConfig bool `yaml:"reports_remote_config"`
}

type Agent struct {
	Executable      string `yaml:"executable"`
	HealthCheckPort int    `yaml:"health_check_port"`
}

type Storage struct {
	Directory string `yaml:"directory"`
}

type SupervisorConfig struct {
	Server       Server       `yaml:"server"`
	Capabilities Capabilities `yaml:"capabilities"`
	Agent        Agent        `yaml:"agent"`
	Storage      Storage      `yaml:"storage"`
}

func initializeSupervisorConfig(logger *zap.Logger, installDir, backupDir string, rb rollback.Rollbacker, hcePort int) error {
	// Open manager.yaml
	managerPath := filepath.Join(installDir, "manager.yaml")
	data, err := os.ReadFile(managerPath)
	if err != nil {
		return fmt.Errorf("read manager.yaml: %w", err)
	}
	var manager map[string]any
	err = yaml.Unmarshal(data, &manager)
	if err != nil {
		return fmt.Errorf("unmarshal manager yaml: %w", err)
	}

	supervisorCfg, err := createSupervisorConfig(logger, manager, installDir, hcePort)
	if err != nil {
		return fmt.Errorf("create supervisor config: %w", err)
	}

	err = writeSupervisorConfig(logger, supervisorCfg, installDir, managerPath, rb)
	if err != nil {
		return fmt.Errorf("write supervisor config")
	}

	err = handleAgentIDConversion(logger, manager, installDir, rb)
	if err != nil {
		return fmt.Errorf("convert agent id: %w", err)
	}

	return nil
}

func createSupervisorConfig(logger *zap.Logger, manager map[string]any, installDir string, hcePort int) (*SupervisorConfig, error) {
	// Read manager values
	var ok bool
	var endpoint, secretKey string
	if endpoint, ok = manager["endpoint"].(string); !ok {
		return nil, errors.New("read in endpoint")
	}
	if secretKey, ok = manager["secret_key"].(string); !ok {
		return nil, errors.New("read in secret key")
	}

	// Construct new supervisor config
	return &SupervisorConfig{
		Server: Server{
			Endpoint: endpoint,
			Headers: Headers{
				Authorization: "Secret-Key " + secretKey,
				AgentIDFormat: "ULID",
			},
			TLS: TLS{
				Insecure:           true,
				InsecureSkipVerify: true,
			},
		},
		Capabilities: Capabilities{
			AcceptsRemoteConfig: true,
			ReportsRemoteConfig: true,
		},
		Agent: Agent{
			Executable:      filepath.Join(installDir, "observiq-otel-collector"),
			HealthCheckPort: hcePort,
		},
		Storage: Storage{
			Directory: filepath.Join(installDir, "supervisor_storage"),
		},
	}, nil
}

func writeSupervisorConfig(logger *zap.Logger, supervisorCfg *SupervisorConfig, installDir, managerPath string, rb rollback.Rollbacker) error {
	supervisorYaml, err := yaml.Marshal(supervisorCfg)
	if err != nil {
		return fmt.Errorf("marshal supervisor yaml: %w", err)
	}

	supervisorPath := filepath.Join(installDir, "supervisor.yaml")

	// We create the action record here, because we want to record the file does not exist
	// before we open the file (which will end up creating the file).
	// Use CopyFileAction because rollback can be used to delete supervisor.yaml.
	// Leave fromPathRel & backupDir empty bc supervisor is a new file.
	cfa, err := action.NewCopyFileAction(logger, "", supervisorPath, "")
	if err != nil {
		return fmt.Errorf("create supervisor.yaml copy file action: %w", err)
	}
	rb.AppendAction(cfa)

	supervisorFile, err := os.OpenFile(filepath.Clean(supervisorPath), os.O_CREATE|os.O_WRONLY, fs.FileMode(0600))
	if err != nil {
		return fmt.Errorf("open supervisor config file: %w", err)
	}
	defer func() {
		err := supervisorFile.Close()
		if err != nil {
			logger.Error("close supervisor config file", zap.Error(err))
		}
	}()

	_, err = supervisorFile.Write(supervisorYaml)
	if err != nil {
		return fmt.Errorf("write supervisor config: %w", err)
	}
	return nil
}

type PersistentState struct {
	InstanceID string `yaml:"instance_id"`
}

func handleAgentIDConversion(logger *zap.Logger, manager map[string]any, installDir string, rb rollback.Rollbacker) error {
	// Retrieve agent id and convert to UUID
	agentID, ok := manager["agent_id"].(string)
	if !ok {
		return fmt.Errorf("read in agent id")
	}
	agentULID, err := ulid.Parse(agentID)
	if err != nil {
		return fmt.Errorf("parse agent id into a ULID: %w", err)
	}
	ulidBytes, err := agentULID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal ulid to bytes: %w", err)
	}
	agentUUID, err := uuid.FromBytes(ulidBytes)
	if err != nil {
		return fmt.Errorf("convert agent id to UUID: %w", err)
	}
	persistentStateYaml, err := yaml.Marshal(&PersistentState{InstanceID: agentUUID.String()})
	if err != nil {
		return fmt.Errorf("marshal persistent_state yaml: %w", err)
	}

	// Create supervisor_storage dir and persistent_state.yaml
	storageDirPath := filepath.Join(installDir, "supervisor_storage")
	persistentStatePath := filepath.Join(storageDirPath, "persistent_state.yaml")

	// Create rollback action for creating 'supervisor_storage' dir.
	// If update fails, rollbacker will delete this dir && persistent_state.yaml file inside it.
	cfa, err := action.NewCopyFileAction(logger, "", storageDirPath, "")
	if err != nil {
		return fmt.Errorf("create supervisor_storage copy file action")
	}
	rb.AppendAction(cfa)

	if err := os.Mkdir(storageDirPath, 0750); err != nil {
		return fmt.Errorf("create supervisor_storage dir: %w", err)
	}

	persistentStateFile, err := os.OpenFile(filepath.Clean(persistentStatePath), os.O_CREATE|os.O_WRONLY, fs.FileMode(0600))
	if err != nil {
		return fmt.Errorf("open persistent_state.yaml: %w", err)
	}
	defer func() {
		err := persistentStateFile.Close()
		if err != nil {
			logger.Error("close persistent_state.yaml file", zap.Error(err))
		}
	}()

	_, err = persistentStateFile.Write(persistentStateYaml)
	if err != nil {
		return fmt.Errorf("write persistent_state.yaml: %w", err)
	}
	return nil
}

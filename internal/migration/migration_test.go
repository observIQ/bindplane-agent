package migration

import (
	"os"
	"path"
	"testing"

	"github.com/observiq/observiq-collector/internal/env"
	"github.com/observiq/observiq-collector/internal/logging"
	"github.com/stretchr/testify/require"
)

var testBaseDir = path.Join(".", "tmp")
var bpAgentConfigDir = path.Join(testBaseDir, "bpagent")
var collectorConfigDir = path.Join(testBaseDir, "collector")

func setupTestDir() error {
	err := os.MkdirAll(testBaseDir, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(bpAgentConfigDir, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(collectorConfigDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func tearDownTestDir() error {
	return os.RemoveAll(testBaseDir)
}

type testBpEnvProvider struct {
	loggingConfigPath string
	remoteConfigPath  string
}

var _ env.BPEnvProvider = env.BPEnvProvider(testBpEnvProvider{})

func (p testBpEnvProvider) RemoteConfig() string {
	return p.remoteConfigPath
}

func (p testBpEnvProvider) LoggingConfig() string {
	return p.loggingConfigPath
}

type testEnvProvider struct {
	loggingConfigPath string
	managerConfigPath string
}

var _ env.EnvProvider = env.EnvProvider(testEnvProvider{})

func (testEnvProvider) LogDir() string {
	return path.Join(testBaseDir, "log")
}

func (p testEnvProvider) DefaultManagerConfigFile() string {
	return p.managerConfigPath
}

func (p testEnvProvider) DefaultLoggingConfigFile() string {
	return p.loggingConfigPath
}

func TestShouldMigrateNoBPAgent(t *testing.T) {
	err := setupTestDir()
	require.NoError(t, err)

	bpEnvProvider := testBpEnvProvider{
		loggingConfigPath: path.Join(bpAgentConfigDir, "does_not_exist.yaml"),
		remoteConfigPath:  path.Join(bpAgentConfigDir, "does_not_exist.yaml"),
	}

	migrate, err := ShouldMigrate(bpEnvProvider)

	require.NoError(t, err)
	require.False(t, migrate)

	err = tearDownTestDir()
	require.NoError(t, err)
}

func TestMigrateMain(t *testing.T) {
	err := setupTestDir()
	require.NoError(t, err)

	bpEnvProvider := testBpEnvProvider{
		loggingConfigPath: path.Join("testdata", "migration", "main_logging.yaml"),
		remoteConfigPath:  path.Join("testdata", "migration", "main_remote.yaml"),
	}

	migrate, err := ShouldMigrate(bpEnvProvider)

	require.NoError(t, err)
	require.True(t, migrate)

	envProvider := testEnvProvider{
		loggingConfigPath: path.Join(collectorConfigDir, "collector-logging.yaml"),
		managerConfigPath: path.Join(collectorConfigDir, "collector-manager.yaml"),
	}

	err = Migrate(envProvider, bpEnvProvider)
	require.NoError(t, err)

	logConf, err := logging.LoadConfig(envProvider.loggingConfigPath)
	require.NoError(t, err)

	// defaultLoggingConf, err := logging.DefaultConfig()
	// require.NoError(t, err)

	// require.Equal(t, &logging.Config{
	// 	Collector: logging.LoggerConfig{
	// 		Level:        zap.InfoLevel,
	// 		MaxBackups:   5,
	// 		MaxMegabytes: 1,
	// 		MaxDays:      7,
	// 		File:         defaultLoggingConf.Collector.File,
	// 	},
	// 	Manager: logging.LoggerConfig{
	// 		Level:        zap.InfoLevel,
	// 		MaxBackups:   5,
	// 		MaxMegabytes: 1,
	// 		MaxDays:      7,
	// 		File:         defaultLoggingConf.Manager.File,
	// 	},
	// }, logConf)

	require.NotNil(t, logConf)

	err = tearDownTestDir()
	require.NoError(t, err)
}

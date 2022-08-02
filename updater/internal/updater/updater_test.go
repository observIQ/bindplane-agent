package updater

import (
	"testing"

	install_mocks "github.com/observiq/observiq-otel-collector/updater/internal/install/mocks"
	rollback_mocks "github.com/observiq/observiq-otel-collector/updater/internal/rollback/mocks"
	service_mocks "github.com/observiq/observiq-otel-collector/updater/internal/service/mocks"
	state_mocks "github.com/observiq/observiq-otel-collector/updater/internal/state/mocks"
	"go.uber.org/zap"
)

func TestUpdaterUpdate(t *testing.T) {
	t.Run("Update is successful", func(t *testing.T) {
		installDir := t.TempDir()
		updater := &Updater{
			installDir: installDir,
			installer:  install_mocks.NewInstaller(t),
			svc:        service_mocks.NewService(t),
			rollbacker: rollback_mocks.NewRollbacker(t),
			monitor:    state_mocks.NewMockMonitor(t),
			logger:     zap.NewNop(),
		}
	})
}

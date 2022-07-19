package rollback

import (
	"testing"

	"github.com/observiq/observiq-otel-collector/updater/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

func TestServiceStartAction(t *testing.T) {
	svc := mocks.NewService(t)
	ssa := NewServiceStartAction(svc)

	svc.On("Stop").Once().Return(nil)

	err := ssa.Rollback()
	require.NoError(t, err)
}

func TestServiceStopAction(t *testing.T) {
	svc := mocks.NewService(t)
	ssa := NewServiceStopAction(svc)

	svc.On("Start").Once().Return(nil)

	err := ssa.Rollback()
	require.NoError(t, err)
}

func TestServiceUpdateAction(t *testing.T) {
	svc := mocks.NewService(t)
	sua := NewServiceUpdateAction("./testdata")
	sua.backupSvc = svc

	svc.On("Update").Once().Return(nil)

	err := sua.Rollback()
	require.NoError(t, err)
}

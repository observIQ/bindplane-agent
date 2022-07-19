package rollback

import (
	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"github.com/observiq/observiq-otel-collector/updater/internal/service"
)

type ServiceStopAction struct {
	svc service.Service
}

func NewServiceStopAction(svc service.Service) ServiceStopAction {
	return ServiceStopAction{
		svc: svc,
	}
}

func (s ServiceStopAction) Rollback() error {
	return s.svc.Start()
}

type ServiceStartAction struct {
	svc service.Service
}

func NewServiceStartAction(svc service.Service) ServiceStartAction {
	return ServiceStartAction{
		svc: svc,
	}
}

func (s ServiceStartAction) Rollback() error {
	return s.svc.Stop()
}

type ServiceUpdateAction struct {
	backupSvc service.Service
}

func NewServiceUpdateAction(tmpDir string) ServiceUpdateAction {
	return ServiceUpdateAction{
		backupSvc: service.NewService(path.BackupDirFromTempDir(tmpDir)),
	}
}

func (s ServiceUpdateAction) Rollback() error {
	return s.backupSvc.Update()
}

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

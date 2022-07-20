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

package path

import "path/filepath"

const (
	latestDirFragment         = "latest"
	rollbackDirFragment       = "rollback"
	serviceFileDirFragment    = "install"
	serviceFileBackupFilename = "backup.service"
)

func LatestDirFromTempDir(tmpDir string) string {
	return filepath.Join(tmpDir, latestDirFragment)
}

func BackupDirFromTempDir(tmpDir string) string {
	return filepath.Join(tmpDir, rollbackDirFragment)
}

func ServiceFileDir(installBaseDir string) string {
	return filepath.Join(installBaseDir, serviceFileDirFragment)
}

func BackupServiceFile(serviceFileDir string) string {
	return filepath.Join(serviceFileDir, serviceFileBackupFilename)
}

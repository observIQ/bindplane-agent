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

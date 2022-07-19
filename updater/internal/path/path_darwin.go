package path

const DarwinInstallDir = "/opt/observiq-otel-collector"

// InstallDir returns the filepath to the install directory
func InstallDir() (string, error) {
	return DarwinInstallDir, nil
}

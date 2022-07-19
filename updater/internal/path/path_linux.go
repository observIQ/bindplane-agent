package paths\

const LinuxInstallDir = "/opt/observiq-otel-collector"

// InstallDir returns the filepath to the install directory
func InstallDir() (string, error) {
	return LinuxInstallDir, nil
}

package install

import "bytes"

// Service represents a controllable service
type Service interface {
	// Start the service
	Start() error

	// Stop the service
	Stop() error

	// Installs the service
	Install() error

	// Uninstalls the service
	Uninstall() error
}

// This function replaces "[INSTALLDIR]" with the given installDir string.
// This is meant to mimic windows "formatted" string syntax.
func replaceInstallDir(b []byte, installDir string) []byte {
	return bytes.ReplaceAll(b, []byte("[INSTALLDIR]"), []byte(installDir))
}

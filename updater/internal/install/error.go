package install

import "fmt"

// InstallStep is a type that represents the step that an installation failed on.
type InstallStep string

const (
	InstallStepCopyFiles        InstallStep = "copy_files"
	InstallStepUninstallService InstallStep = "uninstall_service"
	InstallStepStopService      InstallStep = "stop_service"
	InstallStepInstallService   InstallStep = "install_service"
	InstallStepStartService     InstallStep = "start_service"
	InstallStepHealthCheck      InstallStep = "health_check"
)

// InstallError is an error that wraps another error, also including the step on which the installation failed.
type InstallError struct {
	Step            InstallStep
	UnderlyingError error
}

var _ error = (*InstallError)(nil)

// Unwrap returns the underlying error
func (i InstallError) Unwrap() error {
	return i.UnderlyingError
}

// Error returns the formatted error string
func (i InstallError) Error() string {
	return fmt.Sprintf("failed at step %s: %s", i.Step, i.UnderlyingError.Error())
}

// NewInstallError returns a new InstallError that failed at the given InstallStep, wrapping the given error.
func NewInstallError(location InstallStep, err error) *InstallError {
	return &InstallError{
		Step:            location,
		UnderlyingError: err,
	}
}

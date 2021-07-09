// +build windows

package main

import (
	"os"

	"go.opentelemetry.io/collector/service"
	"golang.org/x/sys/windows/svc"
)

const FLAG_RUN_INTERACTIVE = "OIQCOL_FORCE_INTERACTIVE"

func interactive() (bool, error) {
	if value, present := os.LookupEnv(FLAG_RUN_INTERACTIVE); present && value != "0" {
		return true, nil
	}

	winService, err := svc.IsWindowsService()
	if err != nil {
		return false, err
	}

	return !winService, nil
}

func run(params service.CollectorSettings) error {
	ri, err := interactive()
	if err != nil {
		return err
	}

	if ri {
		return runInteractive(params)
	}

	// This is currently being run as a windows service -- we should run this as a service then.
	return svc.Run("", service.NewWindowsService(params))
}

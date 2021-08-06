// +build !windows

package main

import "github.com/observiq/observiq-collector/manager"

func run(manager *manager.Manager) error {
	return runInteractive(manager)
}

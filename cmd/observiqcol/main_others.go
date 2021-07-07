// +build !windows

package main

import "go.opentelemetry.io/collector/service"

func run(params service.CollectorSettings) error {
	return runInteractive(params)
}

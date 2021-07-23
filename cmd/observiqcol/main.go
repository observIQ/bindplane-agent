package main

import (
	"log"
	"os"

	"github.com/observIQ/observiq-collector/internal/env"
	"github.com/observIQ/observiq-collector/internal/logging"
	"github.com/observIQ/observiq-collector/internal/version"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

func main() {
	factories, err := components()
	if err != nil {
		log.Fatalf("Failed to build default components: %v", err)
	}

	bi := component.BuildInfo{
		Command:     os.Args[0],
		Description: "observIQ's opentelemetry-collector distribution",
		Version:     version.Version,
	}

	if env.IsFileLoggingEnabled() {
		if fp, ok := env.GetLoggingPath(); ok {
			logging.RedirectConsoleOutput(fp)
		} else {
			panic("Failed to find file path for logs, is OIQ_COLLECTOR_HOME set?")
		}
	}

	params := service.CollectorSettings{
		Factories: factories,
		BuildInfo: bi,
	}

	if err := run(params); err != nil {
		log.Fatal(err)
	}

}

func runInteractive(params service.CollectorSettings) error {
	svc, err := service.New(params)
	if err != nil {
		return err
	}

	err = svc.Run()
	if err != nil {
		return err
	}

	return nil
}

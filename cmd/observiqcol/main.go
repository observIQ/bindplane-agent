package main

import (
	"log"
	"os"

	"github.com/observIQ/observIQ-otel-collector/internal/version"

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
		Description: "observIQ opentelemetry-collector distribution",
		Version:     version.Version,
	}

	params := service.CollectorSettings{Factories: factories, BuildInfo: bi}

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

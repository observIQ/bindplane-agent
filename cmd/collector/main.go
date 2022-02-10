package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/observiq/observiq-collector/internal/version"
	"github.com/observiq/observiq-collector/pkg/collector"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/collector/service"
)

func main() {
	var configPaths = pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	pflag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	settings := collector.NewSettings(*configPaths, version.Version(), nil)
	svc, err := service.New(settings)
	if err != nil {
		log.Panicf("Failed to create service: %s", err)
	}

	err = svc.Run(ctx)
	if err != nil {
		log.Panicf("Service received error: %s", err)
	}
}

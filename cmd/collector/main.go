package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/observiq/observiq-collector/internal/version"
	"github.com/observiq/observiq-collector/pkg/collector"
	"github.com/spf13/pflag"
)

func main() {
	var configPath = pflag.String("config", "./config.yaml", "the collector config path")
	pflag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	collector := collector.New(*configPath, version.Version(), nil)
	if err := collector.Run(); err != nil {
		log.Panicf("Collector failed to start: %s", err)
	}
	defer collector.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			status := collector.Status()
			if !status.Running {
				log.Panicf("Collector stopped running: %s", status.Err)
			}
		}
	}
}

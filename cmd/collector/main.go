// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/version"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

const defaultFileLogLevel = zap.InfoLevel

func main() {
	var configPaths = pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	settings, err := collector.NewSettings(*configPaths, version.Version(), nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(ctx, *settings); err != nil {
		log.Fatal(err)
	}
}

func runInteractive(ctx context.Context, params service.CollectorSettings) error {
	svc, err := service.New(params)
	if err != nil {
		return fmt.Errorf("failed to create new service: %w", err)
	}

	if err := svc.Run(ctx); err != nil {
		return fmt.Errorf("collector server run finished with error: %w", err)
	}

	return nil
}

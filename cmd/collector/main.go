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

	"github.com/observiq/observiq-collector/collector"
	"github.com/observiq/observiq-collector/internal/version"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/collector/service"
)

func main() {
	var configPaths = pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	pflag.Parse()

	// TODO: Add this back in when https://github.com/open-telemetry/opentelemetry-collector/issues/4842 is resolved
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer cancel()

	settings := collector.NewSettings(*configPaths, version.Version(), nil)
	if err := run(context.Background(), settings); err != nil {
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

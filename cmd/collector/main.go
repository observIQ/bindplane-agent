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

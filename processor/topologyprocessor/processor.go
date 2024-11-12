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

package topologyprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/observiq/bindplane-agent/internal/topology"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	orgIDHeader      = "x-bp-orgid"
	accountIDHeader  = "x-bp-accountid"
	configNameHeader = "x-bp-config"
)

type topologyUpdate struct {
	gw         topology.GatewayConfigInfo
	routeTable map[topology.GatewayConfigInfo]time.Time
}

type topologyProcessor struct {
	logger      *zap.Logger
	enabled     bool
	topology    *topology.TopologyState
	interval    time.Duration
	processorID component.ID
	bindplane   component.ID

	startOnce sync.Once
}

// newTopologyProcessor creates a new topology processor
func newTopologyProcessor(logger *zap.Logger, cfg *Config, processorID component.ID) (*topologyProcessor, error) {
	destGw := topology.GatewayConfigInfo{
		ConfigName: cfg.ConfigName,
		AccountID:  cfg.AccountID,
		OrgID:      cfg.OrgID,
	}
	topology, err := topology.NewTopologyState(destGw)
	if err != nil {
		return nil, fmt.Errorf("create topology state: %w", err)
	}

	return &topologyProcessor{
		logger:      logger,
		topology:    topology,
		processorID: processorID,
		interval:    cfg.Interval,
		startOnce:   sync.Once{},
	}, nil
}

func (tp *topologyProcessor) start(ctx context.Context, host component.Host) error {
	var err error
	fmt.Println("\033[34m TP START \033[0m")
	tp.startOnce.Do(func() {
		fmt.Println("\033[34m TP DOING START \033[0m")
		registry, getRegErr := GetTopologyRegistry(host, tp.bindplane)
		if getRegErr != nil {
			err = fmt.Errorf("get topology registry: %w", getRegErr)
			return
		}

		if registry != nil {
			registerErr := registry.RegisterTopologyState(tp.processorID.String(), tp.topology)
			if registerErr != nil {
				err = fmt.Errorf("register topology: %w", registerErr)
				return
			}
		}
	})

	return err
}

func (tp *topologyProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	tp.processTopologyHeaders(ctx)
	return td, nil
}

func (tp *topologyProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	fmt.Println("\033[34m PROCESS LOGS \033[0m")
	tp.processTopologyHeaders(ctx)
	return ld, nil
}

func (tp *topologyProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	tp.processTopologyHeaders(ctx)
	return md, nil
}

func (tp *topologyProcessor) processTopologyHeaders(ctx context.Context) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		gw := topology.GatewayConfigInfo{
			ConfigName: metadata.Get(configNameHeader)[0],
			AccountID:  metadata.Get(accountIDHeader)[0],
			OrgID:      metadata.Get(orgIDHeader)[0],
		}
		tp.topology.UpsertRoute(ctx, gw)
	}
}

func (tp *topologyProcessor) shutdown(_ context.Context) error {
	unregisterProcessor(tp.processorID)
	return nil
}

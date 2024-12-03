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
	"strings"
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
	organizationIDHeader = "X-Bindplane-Organization-ID"
	accountIDHeader      = "X-Bindplane-Account-ID"
	configNameHeader     = "X-Bindplane-Configuration"
)

type topologyUpdate struct {
	gw         topology.GatewayInfo
	routeTable map[topology.GatewayInfo]time.Time
}

type topologyProcessor struct {
	logger      *zap.Logger
	enabled     bool
	topology    *topology.ConfigTopology
	interval    time.Duration
	processorID component.ID
	bindplane   component.ID

	startOnce sync.Once
}

// newTopologyProcessor creates a new topology processor
func newTopologyProcessor(logger *zap.Logger, cfg *Config, processorID component.ID) (*topologyProcessor, error) {
	destGw := topology.GatewayInfo{
		GatewayID:  strings.TrimPrefix(cfg.GatewayID, "bindplane_gateway/"),
		ConfigName: cfg.ConfigName,
		AccountID:  cfg.AccountID,
		OrgID:      cfg.OrgID,
	}
	topology, err := topology.NewConfigTopology(destGw)
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

func (tp *topologyProcessor) start(_ context.Context, host component.Host) error {
	var err error
	tp.startOnce.Do(func() {
		registry, getRegErr := GetTopologyRegistry(host, tp.bindplane)
		if getRegErr != nil {
			err = fmt.Errorf("get topology registry: %w", getRegErr)
			return
		}

		if registry != nil {
			registerErr := registry.RegisterConfigTopology(tp.processorID.String(), tp.topology)
			if registerErr != nil {
				return
			}
			registry.SetIntervalChan() <- tp.interval
		}
	})

	return err
}

func (tp *topologyProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	tp.processTopologyHeaders(ctx)
	return td, nil
}

func (tp *topologyProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
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
		var configName, accountID, orgID string

		configNameHeaders := metadata.Get(configNameHeader)
		if len(configNameHeader) > 0 {
			configName = configNameHeaders[0]
		}

		accountIDHeaders := metadata.Get(accountIDHeader)
		if len(configNameHeader) > 0 {
			accountID = accountIDHeaders[0]
		}

		orgIDHeaders := metadata.Get(organizationIDHeader)
		if len(configNameHeader) > 0 {
			orgID = orgIDHeaders[0]
		}

		// only upsert if all headers are present
		if configName != "" && accountID != "" && orgID != "" {
			gw := topology.GatewayInfo{
				ConfigName: configName,
				AccountID:  accountID,
				OrgID:      orgID,
			}
			tp.topology.UpsertRoute(ctx, gw)
		}
	}
}

func (tp *topologyProcessor) shutdown(_ context.Context) error {
	unregisterProcessor(tp.processorID)
	return nil
}

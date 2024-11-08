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
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/opampcustommessages"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	topologyCapability = "com.bindplane.topology"
	orgIDHeader        = "x-bp-orgid"
	accountIDHeader    = "x-bp-accountid"
	configNameHeader   = "x-bp-config"
)

type gateway struct {
	configName string
	accountID  string
	orgID      string
}

type topologyProcessor struct {
	logger *zap.Logger

	processorID      component.ID
	opampExtensionID component.ID

	routeTable map[gateway][]gateway

	customCapabilityHandler opampcustommessages.CustomCapabilityHandler

	started  *atomic.Bool
	stopped  *atomic.Bool
	doneChan chan struct{}
	wg       *sync.WaitGroup
}

// newTopologyProcessor creates a new topology processor
func newTopologyProcessor(logger *zap.Logger, cfg *Config, processorID component.ID) *topologyProcessor {
	return &topologyProcessor{
		logger: logger,

		processorID:      processorID,
		opampExtensionID: cfg.OpAMP,

		routeTable: make(map[gateway][]gateway),

		started:  &atomic.Bool{},
		stopped:  &atomic.Bool{},
		doneChan: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

func (tp *topologyProcessor) start(_ context.Context, host component.Host) error {
	if tp.started.Swap(true) {
		// Start logic should only be run once
		return nil
	}

	ext, ok := host.GetExtensions()[tp.opampExtensionID]
	if !ok {
		return fmt.Errorf("opamp extension %q does not exist", tp.opampExtensionID)
	}

	registry, ok := ext.(opampcustommessages.CustomCapabilityRegistry)
	if !ok {
		return fmt.Errorf("extension %q is not an custom message registry", tp.opampExtensionID)
	}

	var err error
	tp.customCapabilityHandler, err = registry.Register(topologyCapability)
	if err != nil {
		return fmt.Errorf("register custom capability: %w", err)
	}

	return nil
}

func (tp *topologyProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(metadata)
	}
	return td, nil
}

func (tp *topologyProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(metadata)
	}
	return ld, nil
}

func (tp *topologyProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println(metadata)
	}
	return md, nil
}

func (tp *topologyProcessor) stop(ctx context.Context) error {
	if tp.stopped.Swap(true) {
		// Stop logic should only be run once
		return nil
	}

	unregisterProcessor(tp.processorID)

	if tp.customCapabilityHandler != nil {
		tp.customCapabilityHandler.Unregister()
	}

	if tp.doneChan != nil {
		close(tp.doneChan)
	}

	waitgroupDone := make(chan struct{})
	go func() {
		tp.wg.Wait()
		close(waitgroupDone)
	}()

	select {
	case <-waitgroupDone:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// compress gzip compresses the input data
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package observiq

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/golang/snappy"
	"github.com/observiq/bindplane-agent/internal/topology"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// TopologyReporter represents an object that reports topology state.
type TopologyReporter interface {
	TopologyStates() []topology.TopologyState
}

// topologySender is a struct that handles periodically sending topology state via a custom message to an OpAMP endpoint.
type topologySender struct {
	logger      *zap.Logger
	reporter    TopologyReporter
	opampClient client.OpAMPClient
	interval    time.Duration

	changeIntervalChan   chan time.Duration
	changeAttributesChan chan map[string]string

	mux       *sync.Mutex
	isRunning bool
	done      chan struct{}
	wg        *sync.WaitGroup
}

func newTopologySender(l *zap.Logger, reporter TopologyReporter, opampClient client.OpAMPClient, interval time.Duration) *topologySender {
	return &topologySender{
		logger:      l,
		reporter:    reporter,
		opampClient: opampClient,
		interval:    interval,

		changeIntervalChan: make(chan time.Duration, 1),
		mux:                &sync.Mutex{},
		isRunning:          false,
		done:               make(chan struct{}),
		wg:                 &sync.WaitGroup{},
	}
}

// Start starts the sender. It may be called multiple times, even if the sender is already started.
func (ts *topologySender) Start() {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	if ts.isRunning {
		return
	}

	ts.isRunning = true

	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		ts.loop()
	}()
}

// SetInterval changes the interval of the measurements sender.
func (ts topologySender) SetInterval(d time.Duration) {
	select {
	case ts.changeIntervalChan <- d:
	case <-ts.done:
	}

}

func (ts *topologySender) Stop() {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	if !ts.isRunning {
		return
	}

	close(ts.done)
	ts.wg.Wait()

	ts.isRunning = false
}

func (ts *topologySender) loop() {
	t := newTicker()
	t.SetInterval(ts.interval)
	defer t.Stop()

	for {
		select {
		case newInterval := <-ts.changeIntervalChan:
			ts.interval = newInterval
			t.SetInterval(newInterval)
		case <-ts.done:
			return
		case <-t.Chan():
			if ts.reporter == nil {
				// Continue if no reporter available
				ts.logger.Debug("No reporter, skipping sending topology.")
				continue
			}

			topoState := ts.reporter.TopologyStates()
			if len(topoState) == 0 {
				// don't report empty payloads
				continue
			}

			// Send topology state snappy-encoded
			marshalled, err := json.Marshal(topoState)
			if err != nil {
				ts.logger.Error("Failed to marshal topology state.", zap.Error(err))
				continue
			}

			encoded := snappy.Encode(nil, marshalled)

			cm := &protobufs.CustomMessage{
				Capability: topology.ReportTopologyV1Capability,
				Type:       topology.ReportTopologyType,
				Data:       encoded,
			}

			for i := 0; i < maxSendRetries; i++ {
				sendingChannel, err := ts.opampClient.SendCustomMessage(cm)
				switch {
				case err == nil: // OK
				case errors.Is(err, types.ErrCustomMessagePending):
					if i == maxSendRetries-1 {
						// Bail out early, since we aren't going to try to send again
						ts.logger.Warn("Topology were blocked by other custom messages, skipping...", zap.Int("retries", maxSendRetries))
						break
					}

					select {
					case <-sendingChannel:
						continue
					case <-ts.done:
						return
					}
				default:
					ts.logger.Error("Failed to report topology", zap.Error(err))
				}
				break
			}
		}
	}
}

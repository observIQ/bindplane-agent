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

package logcountprocessor

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var (
	// receivers is a map of registered receivers.
	receivers = map[component.ID]*receiver{}

	// receiverMux is a mutex for accessing receivers.
	receiverMux sync.RWMutex

	// errConsumerNotSet is an error returned when a consumer is not set.
	errConsumerNotSet = fmt.Errorf("consumer not set")

	// errReceiverNotSet is an error returned when a receiver is not set.
	errReceiverNotSet = fmt.Errorf("receiver not set")
)

// receiver is the struct that receives log based metrics.
type receiver struct {
	id       component.ID
	consumer consumer.Metrics
}

// newReceiver creates a new receiver.
func newReceiver(id component.ID, consumer consumer.Metrics) *receiver {
	return &receiver{
		id:       id,
		consumer: consumer,
	}
}

// consume consumes log based metrics.
func (r *receiver) consume(ctx context.Context, md pmetric.Metrics) error {
	if r.consumer == nil {
		return errConsumerNotSet
	}

	return r.consumer.ConsumeMetrics(ctx, md)
}

// Start starts the receiver.
func (r *receiver) Start(_ context.Context, _ component.Host) error {
	receiverMux.Lock()
	defer receiverMux.Unlock()
	receivers[r.id] = r

	return nil
}

// Shutdown stops the receiver.
func (r *receiver) Shutdown(_ context.Context) error {
	receiverMux.Lock()
	defer receiverMux.Unlock()
	delete(receivers, r.id)

	return nil
}

// sendMetrics sends log based metrics to the registered receiver.
func sendMetrics(ctx context.Context, id component.ID, md pmetric.Metrics) error {
	receiverMux.RLock()
	defer receiverMux.RUnlock()

	receiver, ok := receivers[id]
	if !ok {
		return errReceiverNotSet
	}

	return receiver.consume(ctx, md)
}

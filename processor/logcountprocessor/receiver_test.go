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
	"testing"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestReceiver(t *testing.T) {
	testCases := []struct {
		name        string
		id          component.ID
		receiver    *receiver
		expectedErr error
	}{
		{
			name: "missing consumer",
			id:   component.NewID("test"),
			receiver: &receiver{
				id: component.NewID("test"),
			},
			expectedErr: errConsumerNotSet,
		},
		{
			name: "missing receiver",
			id:   component.NewID("different"),
			receiver: &receiver{
				id: component.NewID("test"),
			},
			expectedErr: errReceiverNotSet,
		},
		{
			name: "valid receiver",
			id:   component.NewID("test"),
			receiver: &receiver{
				id:       component.NewID("test"),
				consumer: &nopConsumer{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = tc.receiver.Start(context.Background(), nil)
			defer tc.receiver.Shutdown(context.Background())

			err := sendMetrics(context.Background(), tc.id, pmetric.NewMetrics())
			if err != tc.expectedErr {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

// nopConsumer is a nop consumer.
type nopConsumer struct{}

// ConsumeMetrics implements consumer.Metrics.
func (n *nopConsumer) ConsumeMetrics(_ context.Context, _ pmetric.Metrics) error {
	return nil
}

// Capabilities implements consumer.Capabilities.
func (n *nopConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

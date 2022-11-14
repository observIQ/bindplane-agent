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

package throughputwrapper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

func Test_newTraceConsumer(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	tConsumer := newTraceConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, nopLogger, tConsumer.logger)
	require.Equal(t, baseConsumer, tConsumer.baseConsumer)
	require.Len(t, tConsumer.mutators, 1)
	require.Equal(t, &ptrace.ProtoMarshaler{}, tConsumer.tracesSizer)
}

func Test_traceConsumer_ConsumeTraces(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := new(consumertest.TracesSink)
	tConsumer := newTraceConsumer(nopLogger, componentID, baseConsumer)

	td := ptrace.NewTraces()
	td.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()

	err := tConsumer.ConsumeTraces(context.Background(), td)
	require.NoError(t, err)

	require.Equal(t, 1, baseConsumer.SpanCount())
}

func Test_traceConsumer_Capabilities(t *testing.T) {
	nopLogger := zap.NewNop()
	componentID := "id"
	baseConsumer := consumertest.NewNop()
	tConsumer := newTraceConsumer(nopLogger, componentID, baseConsumer)

	require.Equal(t, baseConsumer.Capabilities(), tConsumer.Capabilities())
}

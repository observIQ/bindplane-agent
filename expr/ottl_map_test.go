package expr

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestExtractAttributes(t *testing.T) {
	attrMap, err := MakeOTTLAttributeMap[ottllog.TransformContext](
		map[string]string{
			"key1":           `body["body_key"]`,
			"attr1":          `attributes["key1"]`,
			"does-not-exist": `attributes["dne"]`,
		},
		componenttest.NewNopTelemetrySettings(),
		NewOTTLLogRecordExpression,
	)
	require.NoError(t, err)

	tCtx := ottllog.NewTransformContext(testLogRecord(t), testScope(t), testResource(t))

	mapOut := attrMap.ExtractAttributes(context.Background(), tCtx)

	require.Equal(t, map[string]any{
		"key1":  "cool-thing",
		"attr1": "val1",
	}, mapOut)

}

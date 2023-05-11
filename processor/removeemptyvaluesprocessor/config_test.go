package removeemptyvaluesprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestValidStruct(t *testing.T) {
	require.NoError(t, componenttest.CheckConfigStruct(&Config{}))
}

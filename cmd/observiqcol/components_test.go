package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComponents(t *testing.T) {
	factories, err := components()

	require.NoError(t, err)
	require.NotNil(t, factories.Exporters["otlp"])
	require.NotNil(t, factories.Exporters["otlphttp"])
	require.NotNil(t, factories.Exporters["logging"])
	require.NotNil(t, factories.Exporters["observiq"])

}

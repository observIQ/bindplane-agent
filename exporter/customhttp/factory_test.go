// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package customhttp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configcompression"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
	ocfg, ok := factory.CreateDefaultConfig().(*Config)
	assert.True(t, ok)
	assert.Equal(t, ocfg.ClientConfig.Endpoint, "")
	assert.Equal(t, ocfg.ClientConfig.Timeout, 30*time.Second, "default timeout is 30 second")
	assert.Equal(t, ocfg.RetryConfig.Enabled, true, "default retry is enabled")
	assert.Equal(t, ocfg.RetryConfig.MaxElapsedTime, 300*time.Second, "default retry MaxElapsedTime")
	assert.Equal(t, ocfg.RetryConfig.InitialInterval, 5*time.Second, "default retry InitialInterval")
	assert.Equal(t, ocfg.RetryConfig.MaxInterval, 30*time.Second, "default retry MaxInterval")
	assert.Equal(t, ocfg.QueueConfig.Enabled, true, "default sending queue is enabled")
	assert.Equal(t, ocfg.Encoding, EncodingProto)
	assert.Equal(t, ocfg.Compression, configcompression.TypeGzip)
}

// func TestCreateLogsExporter(t *testing.T) {
// 	factory := NewFactory()
// 	cfg := factory.CreateDefaultConfig().(*Config)
// 	cfg.ClientConfig.Endpoint = "http://" + testutil.GetAvailableLocalAddress(t)

// 	set := exportertest.NewNopSettings()
// 	oexp, err := factory.CreateLogsExporter(context.Background(), set, cfg)
// 	require.Nil(t, err)
// 	require.NotNil(t, oexp)
// }

func TestComposeSignalURL(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)

	// Has slash at end
	cfg.ClientConfig.Endpoint = "http://localhost:4318/"
	url, err := composeSignalURL(cfg, "", "traces")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:4318/v1/traces", url)

	// No slash at end
	cfg.ClientConfig.Endpoint = "http://localhost:4318"
	url, err = composeSignalURL(cfg, "", "traces")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:4318/v1/traces", url)
}

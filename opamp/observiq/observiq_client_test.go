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

package observiq

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/observiq/observiq-otel-collector/opamp/mocks"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		desc        string
		config      opamp.Config
		expectedErr error
	}{
		{
			desc: "Bad URL Scheme",
			config: opamp.Config{
				Endpoint: "http://localhost:1234",
				AgentID:  "b24181a8-bc16-4ec1-b3af-ca6f7b669af8",
			},
			expectedErr: ErrUnsupportedURL,
		},
		{
			desc: "Invalid Endpoint",
			config: opamp.Config{
				Endpoint: "\t\t\t",
				AgentID:  "b24181a8-bc16-4ec1-b3af-ca6f7b669af8",
			},
			expectedErr: errors.New("net/url: invalid control character in URL"),
		},
		{
			desc: "Valid Config",
			config: opamp.Config{
				Endpoint: "ws://localhost:1234",
				AgentID:  "b24181a8-bc16-4ec1-b3af-ca6f7b669af8",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testLogger := zap.NewNop().Sugar()
			mockManager := new(mocks.MockConfigManager)
			shutdownChan := make(chan struct{})

			actual, err := NewClient(testLogger, tc.config, mockManager, shutdownChan)

			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
				assert.Nil(t, actual)
			} else {
				assert.NoError(t, err)

				observiqClient, ok := actual.(*Client)
				require.True(t, ok)

				// Do a shallow check on all fields to assert they exist and are equal to passed in params were possible
				assert.NotNil(t, observiqClient.opampClient)
				assert.Equal(t, mockManager, observiqClient.configManager)
				assert.Equal(t, shutdownChan, observiqClient.shutdownChan)
				assert.Equal(t, testLogger.Named("opamp"), observiqClient.logger)
				assert.NotNil(t, observiqClient.ident)
			}

		})
	}
}

func TestClientConnect(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "SetAgentDescription fails",
			testFunc: func(*testing.T) {
				expectedErr := errors.New("oops")

				mockOpAmpClient := new(mocks.MockClient)
				mockOpAmpClient.On("SetAgentDescription", mock.Anything).Return(expectedErr)

				c := &Client{
					opampClient:   mockOpAmpClient,
					logger:        zap.NewNop().Sugar(),
					ident:         &identity{},
					configManager: nil,
					shutdownChan:  make(chan struct{}),
				}

				err := c.Connect(context.Background(), opamp.Config{})
				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Start fails",
			testFunc: func(*testing.T) {
				expectedErr := errors.New("oops")

				mockOpAmpClient := new(mocks.MockClient)
				mockOpAmpClient.On("SetAgentDescription", mock.Anything).Return(nil)
				mockOpAmpClient.On("Start", mock.Anything, mock.Anything).Return(expectedErr)

				secretKey := "136bdd08-2074-40b7-ac1c-6706ac24c4f2"
				config := opamp.Config{
					Endpoint:  "ws://localhost:1234",
					SecretKey: &secretKey,
					AgentID:   "a69dcef0-0261-4f4f-9ac0-a483af42a6ba",
				}

				c := &Client{
					opampClient: mockOpAmpClient,
					logger:      zap.NewNop().Sugar(),
					ident: &identity{
						agentID: config.AgentID,
					},
					configManager: nil,
					shutdownChan:  make(chan struct{}),
				}

				err := c.Connect(context.Background(), config)
				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Connect successful",
			testFunc: func(*testing.T) {
				mockOpAmpClient := new(mocks.MockClient)
				mockOpAmpClient.On("SetAgentDescription", mock.Anything).Return(nil)

				secretKey := "136bdd08-2074-40b7-ac1c-6706ac24c4f2"
				config := opamp.Config{
					Endpoint:  "ws://localhost:1234",
					SecretKey: &secretKey,
					AgentID:   "a69dcef0-0261-4f4f-9ac0-a483af42a6ba",
				}

				c := &Client{
					opampClient: mockOpAmpClient,
					logger:      zap.NewNop().Sugar(),
					ident: &identity{
						agentID: config.AgentID,
					},
					configManager: nil,
					shutdownChan:  make(chan struct{}),
				}

				expectedSettings := types.StartSettings{
					OpAMPServerURL:      config.Endpoint,
					AuthorizationHeader: *config.SecretKey,
					TLSConfig:           nil,
					InstanceUid:         c.ident.agentID,
					Callbacks: types.CallbacksStruct{
						OnConnectFunc:          c.onConnectHandler,
						OnConnectFailedFunc:    c.onConnectFailedHandler,
						OnErrorFunc:            c.onErrorHandler,
						OnRemoteConfigFunc:     c.onRemoteConfigHandler,
						GetEffectiveConfigFunc: c.onGetEffectiveConfigHandler,
					},
				}
				mockOpAmpClient.On("Start", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					settings := args.Get(1).(types.StartSettings)
					assert.Equal(t, expectedSettings.OpAMPServerURL, settings.OpAMPServerURL)
					assert.Equal(t, expectedSettings.AuthorizationHeader, settings.AuthorizationHeader)
					assert.Equal(t, expectedSettings.TLSConfig, settings.TLSConfig)
					assert.Equal(t, expectedSettings.InstanceUid, settings.InstanceUid)
					// assert is unable to compare function pointers
				})

				err := c.Connect(context.Background(), config)
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestClientDisconnect(t *testing.T) {
	mockOpAmpClient := new(mocks.MockClient)
	ctx := context.Background()
	mockOpAmpClient.On("Stop", ctx).Return(nil)

	c := &Client{
		opampClient: mockOpAmpClient,
	}

	c.Disconnect(ctx)
	mockOpAmpClient.AssertExpectations(t)
}

func TestClient_onConnectionHandler(t *testing.T) {
	mockOpAmpClient := new(mocks.MockClient)

	c := &Client{
		logger:      zap.NewNop().Sugar(),
		opampClient: mockOpAmpClient,
		ident: &identity{
			agentID:     "4322d8d1-f3e0-46db-b68d-b01a4689ef19",
			agentName:   nil,
			serviceName: "com.observiq.collector",
			version:     "v1.2.3",
			labels:      nil,
			oSArch:      "amd64",
			oSDetails:   "os details",
			oSFamily:    "linux",
			hostname:    "my-linux-box",
			mac:         "68-C7-B4-EB-A8-D2",
		},
	}

	mockOpAmpClient.On("SetAgentDescription", mock.Anything).Return(nil)

	c.onConnectHandler()
	mockOpAmpClient.AssertExpectations(t)
}

func TestClient_onGetEffectiveConfigHandler(t *testing.T) {
	mockManager := new(mocks.MockConfigManager)

	c := &Client{
		logger:        zap.NewNop().Sugar(),
		configManager: mockManager,
	}

	mockManager.On("ComposeEffectiveConfig").Return(&protobufs.EffectiveConfig{}, nil)

	c.onGetEffectiveConfigHandler(context.Background())
	mockManager.AssertExpectations(t)
}

func TestClient_onRemoteConfigHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Config Changes return error",
			testFunc: func(*testing.T) {
				expectedErr := errors.New("oops")
				expectedChanged := false
				mockManager := new(mocks.MockConfigManager)
				mockManager.On("ApplyConfigChanges", mock.Anything).Return(&protobufs.EffectiveConfig{}, expectedChanged, expectedErr)

				c := &Client{
					configManager: mockManager,
					logger:        zap.NewNop().Sugar(),
				}

				effCfg, changed, err := c.onRemoteConfigHandler(context.Background(), &protobufs.AgentRemoteConfig{})
				assert.Nil(t, effCfg)
				assert.Equal(t, expectedChanged, changed)
				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Config Changes occur",
			testFunc: func(*testing.T) {
				expectedEffCfg := &protobufs.EffectiveConfig{}
				mockManager := new(mocks.MockConfigManager)
				mockManager.On("ApplyConfigChanges", mock.Anything).Return(expectedEffCfg, true, nil)

				shutdownChan := make(chan struct{}, 1)
				c := &Client{
					configManager: mockManager,
					logger:        zap.NewNop().Sugar(),
					shutdownChan:  shutdownChan,
				}

				effCfg, changed, err := c.onRemoteConfigHandler(context.Background(), &protobufs.AgentRemoteConfig{})
				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)

				shutDownFunc := func() bool {
					<-shutdownChan
					return true
				}
				assert.Eventually(t, shutDownFunc, 1*time.Minute, 200*time.Millisecond)
			},
		},
		{
			desc: "No Config Changes occur",
			testFunc: func(*testing.T) {
				expectedEffCfg := &protobufs.EffectiveConfig{}
				mockManager := new(mocks.MockConfigManager)
				mockManager.On("ApplyConfigChanges", mock.Anything).Return(expectedEffCfg, false, nil)

				shutdownChan := make(chan struct{}, 1)
				c := &Client{
					configManager: mockManager,
					logger:        zap.NewNop().Sugar(),
					shutdownChan:  shutdownChan,
				}

				effCfg, changed, err := c.onRemoteConfigHandler(context.Background(), &protobufs.AgentRemoteConfig{})
				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)

				// If we pushed to a closed channel we should panic
				// if channel is still open means it didn't signal a shutdown
				testShutdownOpen := func() {
					shutdownChan <- struct{}{}

				}
				assert.NotPanics(t, testShutdownOpen)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

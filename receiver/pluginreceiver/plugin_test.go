package pluginreceiver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
)

func TestLoadPlugin(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		expectedResult *Plugin
		expectedErr    error
	}{
		{
			name: "valid plugin",
			path: "./test/plugin-valid.yaml",
			expectedResult: &Plugin{
				Title:       "test-plugin",
				Template:    "receivers:",
				Version:     "0.0.0",
				Description: "A valid test plugin",
				Parameters: []Parameter{
					{
						Name:      "env",
						Type:      stringType,
						Default:   "prod",
						Supported: []interface{}{"prod", "dev"},
						Required:  true,
					},
				},
			},
		},
		{
			name:        "invalid plugin",
			path:        "./test/plugin-invalid-yaml.yaml",
			expectedErr: errors.New("failed to unmarshal plugin from yaml"),
		},
		{
			name:        "missing file",
			path:        "./test/missing.yaml",
			expectedErr: errors.New("failed to read file"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := LoadPlugin(tc.path)
			switch tc.expectedErr {
			case nil:
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			default:
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestRenderConfig(t *testing.T) {
	testCases := []struct {
		name           string
		plugin         *Plugin
		values         map[string]interface{}
		expectedResult *config.Map
		expectedErr    error
	}{
		{
			name: "invalid template error",
			plugin: &Plugin{
				Template: "{{.invalid",
			},
			expectedErr: errors.New("failed to create plugin template"),
		},
		{
			name: "template execution error",
			plugin: &Plugin{
				Template: `{{template "base" .}}`,
			},
			expectedErr: errors.New("failed to execute template"),
		},
		{
			name: "invalid yaml error",
			plugin: &Plugin{
				Template: "test template",
			},
			expectedErr: errors.New("failed to unmarshal template result as yaml"),
		},
		{
			name: "valid template",
			plugin: &Plugin{
				Template: "{{if .enabled}}receivers:{{end}}",
			},
			values: map[string]interface{}{
				"enabled": true,
			},
			expectedResult: config.NewMapFromStringMap(map[string]interface{}{
				"receivers": nil,
			}),
		},
		{
			name: "valid template with defaults",
			plugin: &Plugin{
				Template: "{{if .enabled}}receivers:{{end}}",
				Parameters: []Parameter{
					{
						Name:    "enabled",
						Default: true,
					},
				},
			},
			expectedResult: config.NewMapFromStringMap(map[string]interface{}{
				"receivers": nil,
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.plugin.RenderConfig(tc.values)
			switch tc.expectedErr {
			case nil:
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			default:
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	testCases := []struct {
		name           string
		plugin         *Plugin
		values         map[string]interface{}
		expectedResult map[string]interface{}
	}{
		{
			name: "with no parameters",
			plugin: &Plugin{
				Parameters: nil,
			},
			values: map[string]interface{}{
				"param1": "value",
			},
			expectedResult: map[string]interface{}{
				"param1": "value",
			},
		},
		{
			name: "with no defaults",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name: "param1",
					},
					{
						Name: "param2",
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value",
			},
			expectedResult: map[string]interface{}{
				"param1": "value",
			},
		},
		{
			name: "with defaults",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name:    "param1",
						Default: "defaultValue1",
					},
					{
						Name:    "param2",
						Default: "defaultValue2",
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value",
			},
			expectedResult: map[string]interface{}{
				"param1": "value",
				"param2": "defaultValue2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.plugin.ApplyDefaults(tc.values)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCheckParameters(t *testing.T) {
	testCases := []struct {
		name        string
		plugin      *Plugin
		values      map[string]interface{}
		expectedErr error
	}{
		{
			name:   "undefined parameters",
			plugin: &Plugin{},
			values: map[string]interface{}{
				"param1": "value1",
			},
			expectedErr: errors.New("definition failure"),
		},
		{
			name: "missing required",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name:     "param1",
						Required: true,
					},
				},
			},
			values:      map[string]interface{}{},
			expectedErr: errors.New("required failure"),
		},
		{
			name: "invalid string type",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name: "param1",
						Type: stringType,
					},
				},
			},
			values: map[string]interface{}{
				"param1": 5,
			},
			expectedErr: errors.New("must be a string"),
		},
		{
			name: "invalid []string type",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name: "param1",
						Type: stringArrayType,
					},
				},
			},
			values: map[string]interface{}{
				"param1": 5,
			},
			expectedErr: errors.New("must be a []string"),
		},
		{
			name: "invalid int type",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name: "param1",
						Type: intType,
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value1",
			},
			expectedErr: errors.New("must be an int"),
		},
		{
			name: "invalid bool type",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name: "param1",
						Type: boolType,
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value1",
			},
			expectedErr: errors.New("must be a bool"),
		},
		{
			name: "not supported value",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name:      "param1",
						Type:      stringType,
						Supported: []interface{}{"value2"},
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value1",
			},
			expectedErr: errors.New("supported value failure"),
		},
		{
			name: "valid parameters",
			plugin: &Plugin{
				Parameters: []Parameter{
					{
						Name:      "param1",
						Type:      stringType,
						Required:  true,
						Supported: []interface{}{"value1"},
					},
					{
						Name: "param2",
						Type: stringArrayType,
					},
					{
						Name: "param3",
						Type: intType,
					},
					{
						Name: "param4",
						Type: boolType,
					},
					{
						Name:      "param5",
						Type:      stringType,
						Supported: []interface{}{"value5"},
					},
				},
			},
			values: map[string]interface{}{
				"param1": "value1",
				"param2": []string{"value2"},
				"param3": 5,
				"param4": true,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.plugin.CheckParameters(tc.values)
			switch tc.expectedErr {
			case nil:
				require.NoError(t, err)
			default:
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

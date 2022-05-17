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

package opamp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoopValidator(t *testing.T) {
	require.True(t, NoopValidator(nil))
}

func TestNewYamlValidator(t *testing.T) {
	ymlCfg := `
key: value
num: 1
`
	target := make(map[string]interface{})
	validatorFunc := NewYamlValidator(target)

	require.True(t, validatorFunc([]byte(ymlCfg)))
	require.False(t, validatorFunc([]byte("garbage")))
}

func TestNewJSONValidator(t *testing.T) {
	jsonCfg := `{"key":"value"}`
	target := make(map[string]interface{})
	validatorFunc := NewJSONValidator(target)

	require.True(t, validatorFunc([]byte(jsonCfg)))
	require.False(t, validatorFunc([]byte("garbage")))
}

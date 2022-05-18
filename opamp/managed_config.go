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

// ReloadFunc is a function that handles reloading a config given the new contents
type ReloadFunc func([]byte) (changed bool, err error)

// NoopReloadFunc used as a noop reload function if unsure of how to reload
func NoopReloadFunc([]byte) (bool, error) {
	return false, nil
}

// ManagedConfig is a structure that can manager an on disk config file
type ManagedConfig struct {
	// ConfigPath is the path on disk where the configuration lives
	ConfigPath string

	// Reload will be called when any changes to this config occur.
	Reload ReloadFunc
}

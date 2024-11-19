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

// Package counter contains structs used to count telemetry grouped by resource and attributes.
package counter

import "encoding/json"

// TelemetryCounter tracks the number of times a set of resource and attribute dimensions have been seen.
type TelemetryCounter struct {
	resources map[string]*ResourceCounter
	commands  chan func()
}

// NewTelemetryCounter creates a new TelemetryCounter.
func NewTelemetryCounter() *TelemetryCounter {
	t := &TelemetryCounter{
		resources: make(map[string]*ResourceCounter),
		commands:  make(chan func()),
	}
	go t.run()
	return t
}

// run listens for commands to modify or read the resources.
func (t *TelemetryCounter) run() {
	for cmd := range t.commands {
		cmd()
	}
}

// Add increments the counter with the supplied dimensions.
func (t *TelemetryCounter) Add(resource, attributes map[string]any) {
	t.commands <- func() {
		key := getDimensionKey(resource)
		if _, ok := t.resources[key]; !ok {
			t.resources[key] = newResourceCounter(resource)
		}
		t.resources[key].add(attributes)
	}
}

// Resources returns a map of resource ID to a counter for that resource and resets the counter.
func (t *TelemetryCounter) Resources() map[string]*ResourceCounter {
	result := make(chan map[string]*ResourceCounter)
	t.commands <- func() {
		result <- t.resources
		t.resources = make(map[string]*ResourceCounter) // Reset the counter
	}
	return <-result
}

// ResourceCounter dimensions the counter by resource.
type ResourceCounter struct {
	values     map[string]any
	attributes map[string]*AttributeCounter
}

// newResourceCounter creates a new ResourceCounter.
func newResourceCounter(values map[string]any) *ResourceCounter {
	return &ResourceCounter{
		values:     values,
		attributes: map[string]*AttributeCounter{},
	}
}

// add increments the counter with the supplied dimensions.
func (r *ResourceCounter) add(attributes map[string]any) {
	key := getDimensionKey(attributes)
	if _, ok := r.attributes[key]; !ok {
		r.attributes[key] = newAttributeCounter(attributes)
	}

	r.attributes[key].add()
}

// Attributes returns a map of attribute set ID to a counter for that attribute set.
func (r ResourceCounter) Attributes() map[string]*AttributeCounter {
	return r.attributes
}

// Values returns the raw map value of the resource that this counter counts.
func (r ResourceCounter) Values() map[string]any {
	return r.values
}

// AttributeCounter dimensions the counter by attributes.
type AttributeCounter struct {
	values map[string]any
	count  int
}

// newAttributeCounter creates a new AttributeCounter.
func newAttributeCounter(values map[string]any) *AttributeCounter {
	return &AttributeCounter{
		values: values,
	}
}

// add increments the counter.
func (a *AttributeCounter) add() {
	a.count++
}

// Count returns the number of counts for this attribute counter.
func (a AttributeCounter) Count() int {
	return a.count
}

// Values returns the attribute map that this counter tracks.
func (a AttributeCounter) Values() map[string]any {
	return a.values
}

// getDimensionKey returns a unique key for the dimension.
func getDimensionKey(dimension map[string]any) string {
	dimensionJSON, _ := json.Marshal(dimension)
	return string(dimensionJSON)
}

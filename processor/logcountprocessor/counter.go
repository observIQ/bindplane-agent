package logcountprocessor

import "encoding/json"

// LogCounter tracks the number of times a set of resource and attribute dimensions have been seen.
type LogCounter struct {
	resources map[string]*ResourceCounter
}

// NewLogCounter creates a new LogCounter.
func NewLogCounter() *LogCounter {
	return &LogCounter{
		resources: make(map[string]*ResourceCounter),
	}
}

// Add increments the counter with the supplied dimensions.
func (l *LogCounter) Add(resource, attributes map[string]interface{}) {
	key := getDimensionKey(resource)
	if _, ok := l.resources[key]; !ok {
		l.resources[key] = NewResourceCounter(resource)
	}

	l.resources[key].Add(attributes)
}

// Reset resets the counter.
func (l *LogCounter) Reset() {
	l.resources = make(map[string]*ResourceCounter)
}

// ResourceCounter dimensions the counter by resource.
type ResourceCounter struct {
	values     map[string]interface{}
	attributes map[string]*AttributeCounter
}

// NewResourceCounter creates a new ResourceCounter.
func NewResourceCounter(values map[string]interface{}) *ResourceCounter {
	return &ResourceCounter{
		values:     values,
		attributes: map[string]*AttributeCounter{},
	}
}

// Add increments the counter with the supplied dimensions.
func (r *ResourceCounter) Add(attributes map[string]interface{}) {
	key := getDimensionKey(attributes)
	if _, ok := r.attributes[key]; !ok {
		r.attributes[key] = NewAttributeCounter(attributes)
	}

	r.attributes[key].Add()
}

// AttributeCounter dimensions the counter by attributes.
type AttributeCounter struct {
	values map[string]interface{}
	count  int
}

// NewAttributeCounter creates a new AttributeCounter.
func NewAttributeCounter(values map[string]interface{}) *AttributeCounter {
	return &AttributeCounter{
		values: values,
	}
}

// Add increments the counter.
func (a *AttributeCounter) Add() {
	a.count++
}

// getDimensionKey returns a unique key for the dimension.
func getDimensionKey(dimension map[string]interface{}) string {
	dimensionJSON, _ := json.Marshal(dimension)
	return string(dimensionJSON)
}

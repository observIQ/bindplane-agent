package alternateprocessor

import (
	"fmt"
	"strings"
	"time"
)

// timeUnitMap is a map of time units to their associated value
var timeUnitMap = map[string]time.Duration{
	"s":      time.Second,
	"sec":    time.Second,
	"second": time.Second,
	"m":      time.Minute,
	"min":    time.Minute,
	"minute": time.Minute,
	"h":      time.Hour,
	"hour":   time.Hour,
	"d":      24 * time.Hour,
	"day":    24 * time.Hour,
}

// measureUnitMap is a map of measure units to their associated value
var measureUnitMap = map[string]float64{
	"logs":    1,
	"metrics": 1,
	"spans":   1,
	"records": 1,
	"items":   1,
	"b":       1,
	"by":      1,
	"bytes":   1,
	"kib":     1024,
	"kb":      1000,
	"k":       1000,
	"mib":     1024 * 1024,
	"miby":    1024 * 1024,
	"mb":      1000 * 1000,
	"mby":     1000 * 1000,
	"m":       1000 * 1000,
	"gib":     1024 * 1024 * 1024,
	"giby":    1024 * 1024 * 1024,
	"gb":      1000 * 1000 * 1000,
	"g":       1000 * 1000 * 1000,
	"tib":     1024 * 1024 * 1024 * 1024,
	"tb":      1000 * 1000 * 1000 * 1000,
	"t":       1000 * 1000 * 1000 * 1000,
	"pib":     1024 * 1024 * 1024 * 1024 * 1024,
	"pb":      1000 * 1000 * 1000 * 1000 * 1000,
}

// TimeUnit is a unit of time
type TimeUnit struct {
	Name  string
	Value time.Duration
}

// MeasureUnit is a unit of measure
type MeasureUnit struct {
	Name  string
	Value float64
}

// IsLogCount returns true if the measure is a log count
func (m *MeasureUnit) IsLogCount() bool {
	return m.Name == "logs"
}

// IsMetricCount returns true if the measure is a metric count
func (m *MeasureUnit) IsMetricCount() bool {
	return m.Name == "metrics"
}

// IsSpanCount returns true if the measure is a span count
func (m *MeasureUnit) IsSpanCount() bool {
	return m.Name == "spans"
}

// IsTotalCount returns true if the measure is a total count
func (m *MeasureUnit) IsTotalCount() bool {
	return m.Name == "records" || m.Name == "items"
}

// IsSizeCount returns true if the measure is a size count
func (m *MeasureUnit) IsSizeCount() bool {
	switch m.Name {
	case "b", "bytes", "kib", "kb", "k", "mib", "mb", "m", "gib", "gb", "g", "tib", "tb", "t", "pib", "pb":
		return true
	default:
		return false
	}
}

// NewTimeUnit creates a new time unit
func NewTimeUnit(name string, value time.Duration) *TimeUnit {
	return &TimeUnit{
		Name:  name,
		Value: value,
	}
}

// NewMeasureUnit creates a new measure unit
func NewMeasureUnit(name string, value float64) *MeasureUnit {
	return &MeasureUnit{
		Name:  name,
		Value: value,
	}
}

// ParseTimeUnit parses a time unit from a string
func ParseTimeUnit(unit string) (*TimeUnit, error) {
	duration, ok := timeUnitMap[strings.ToLower(unit)]
	if !ok {
		return nil, fmt.Errorf("invalid time unit: %s", unit)
	}

	return NewTimeUnit(unit, duration), nil
}

// ParseMeasureUnit parses a measure unit from a string
func ParseMeasureUnit(unit string) (*MeasureUnit, error) {
	size, ok := measureUnitMap[strings.ToLower(unit)]
	if !ok {
		return nil, fmt.Errorf("invalid measure unit: %s", unit)
	}

	return NewMeasureUnit(unit, size), nil
}

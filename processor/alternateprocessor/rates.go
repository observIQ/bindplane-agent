package alternateprocessor

import (
	"fmt"
	"regexp"
	"strconv"
)

// rateRegex is the regex to match a rate
var rateRegex = regexp.MustCompile(`^(\d+)\s?([a-zA-Z]+)/([a-zA-Z]+)$`)

// Rate is a rate of throughput
type Rate struct {
	// Value is the value of the rate
	Value float64
	// Measure is the measure unit of the rate
	Measure *MeasureUnit
	// Time is the time unit of the rate
	Time *TimeUnit
}

// NormalizedValue returns the normalized value of a rate
func (r *Rate) NormalizedValue() float64 {
	return r.Value * r.Measure.Value / r.Time.Value.Seconds()
}

// String returns the string representation of a rate
func (r *Rate) String() string {
	return fmt.Sprintf("%f%s/%s", r.Value, r.Measure.Name, r.Time.Name)
}

// NewRate creates a new rate
func NewRate(value float64, measure *MeasureUnit, time *TimeUnit) *Rate {
	return &Rate{
		Value:   value,
		Measure: measure,
		Time:    time,
	}
}

// ParseRate parses a rate from a string
func ParseRate(str string) (*Rate, error) {
	matches := rateRegex.FindStringSubmatch(str)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid rate structure: %s", str)
	}
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, fmt.Errorf("rate value is not a number: %s", matches[1])
	}
	measure, err := ParseMeasureUnit(matches[2])
	if err != nil {
		return nil, err
	}
	time, err := ParseTimeUnit(matches[3])
	if err != nil {
		return nil, err
	}
	return NewRate(value, measure, time), nil
}

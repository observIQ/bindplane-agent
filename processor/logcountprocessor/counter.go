package logcountprocessor

import (
	"encoding/json"

	"go.uber.org/zap"
)

// Counter is a counter used to count the instances of log records.
type Counter struct {
	counts map[string]*Count
	logger *zap.Logger
}

// NewCounter returns a new counter.
func NewCounter(logger *zap.Logger) *Counter {
	return &Counter{
		counts: make(map[string]*Count),
		logger: logger,
	}
}

// Add will add a log record to the counter.
func (c *Counter) Add(record Record) {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		c.logger.Error("Failed to add log record to counter", zap.Error(err))
		return
	}
	recordKey := string(recordJSON)

	if _, ok := c.counts[recordKey]; !ok {
		c.counts[recordKey] = NewCount(record)
	}

	c.counts[recordKey].Increment()
}

// Reset resets the record counter.
func (c *Counter) Reset() {
	for k := range c.counts {
		delete(c.counts, k)
	}
}

// Count is the count of a log record.
type Count struct {
	record Record
	value  int
}

// NewCount creates a new log record count.
func NewCount(record Record) *Count {
	return &Count{
		record: record,
	}
}

// Increment increments the count.
func (c *Count) Increment() {
	c.value++
}

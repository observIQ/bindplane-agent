package cloudwatch

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/open-telemetry/opentelemetry-log-collection/operator"
)

type Persister struct {
	DB operator.Persister
}

// Helper function to get persisted data
func (p *Persister) Read(ctx context.Context, key string) (int64, error) {
	startTimeBytes, err := p.DB.Get(ctx, key)
	if err != nil {
		return -1, fmt.Errorf("there was an error reading from persistent storage: %w", err)
	}
	buffer := bytes.NewBuffer(startTimeBytes)
	var startTime int64
	err = binary.Read(buffer, binary.BigEndian, &startTime)
	if err != nil && err.Error() != "EOF" {
		return 0, err
	}
	return startTime, nil
}

// Helper function to set persisted data
func (p *Persister) Write(ctx context.Context, key string, value int64) error {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return p.DB.Set(ctx, key, buf)
}

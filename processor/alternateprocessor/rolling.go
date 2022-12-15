package alternateprocessor

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	errInvalidBucketCount = errors.New("not a valid bucket count")
)

// RollingAverage is a rolling average
type RollingAverage struct {
	// buckets is a circular buffer of values observed in each bucket
	buckets []float64

	// interval is the interval before rolling over to the next bucket
	interval time.Duration

	// cancel is the function to call to cancel the rolling average
	cancel context.CancelFunc

	// wg is the wait group to use to wait for the rolling average to stop
	wg sync.WaitGroup

	// mux is the mutex to use to protect the buckets
	mux *sync.Mutex
}

// AddBytes adds a value to the rolling average
func (r *RollingAverage) AddBytes(value float64) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.buckets[0] += value
}

// Start starts the rolling average
func (r *RollingAverage) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	go r.handleRoll(ctx)
}

// handleRoll is a helper function to handle the rolling of buckets
func (r *RollingAverage) handleRoll(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(r.interval)
	for {
		select {
		case <-ticker.C:
			r.roll()
		case <-ctx.Done():
			return
		}
	}
}

// roll rolls the buckets
func (r *RollingAverage) roll() {
	r.mux.Lock()
	defer r.mux.Unlock()

	for i := len(r.buckets) - 1; i > 0; i-- {
		r.buckets[i] = r.buckets[i-1]
	}
	r.buckets[0] = 0
}

// Stop stops the rolling average
func (r *RollingAverage) Stop() {
	r.cancel()
	r.wg.Wait()
}

// Value returns the current value of the rolling average
func (r *RollingAverage) Value() float64 {
	r.mux.Lock()
	defer r.mux.Unlock()

	var total float64
	for _, bucket := range r.buckets[1:] {
		total += bucket
	}

	bucketCount := len(r.buckets) - 1
	return total / float64(bucketCount)
}

// NormalizedRateValue returns the current value of the rolling average normalized by the interval
func (r *RollingAverage) NormalizedRateValue() float64 {
	return r.Value() / r.interval.Seconds()
}

// NewRollingAverage creates a new rolling average
func NewRollingAverage(bucketCount int, interval time.Duration) (*RollingAverage, error) {
	if bucketCount < 1 {
		return nil, errInvalidBucketCount
	}

	return &RollingAverage{
		mux:      &sync.Mutex{},
		buckets:  make([]float64, bucketCount+1),
		interval: interval,
	}, nil
}

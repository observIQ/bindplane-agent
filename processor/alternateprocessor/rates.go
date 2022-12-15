package alternateprocessor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// secondUnit is the unit of time for a second
	secondUnit TimeUnit = "second"

	// minuteUnit is the unit of time for a minute
	minuteUnit TimeUnit = "minute"

	// hourUnit is the unit of time for an hour
	hourUnit TimeUnit = "hour"

	// dayUnit is the unit of time for a day
	dayUnit TimeUnit = "day"
)

// timeUnitMap is a map of time units to time durations
var timeUnitMap = map[TimeUnit]time.Duration{
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,
	"day":    24 * time.Hour,
}

// errWindowTooSmall is the error returned when the window is smaller than the bucket duration
var errWindowTooSmall = errors.New("window must be greater than or equal to bucket duration")

// TimeUnit is the unit of time to use for the rate
type TimeUnit string

// RateTrackerConfig is the configuration for a rate tracker
type RateTrackerConfig struct {
	// TimeWindow is the window of time to track
	TimeWindow time.Duration `mapstructure:"window"`

	// TimeUnit is the unit of time to use for the rate
	TimeUnit TimeUnit `mapstructure:"unit"`
}

// Build builds a rate tracker from a rate tracker config
func (c *RateTrackerConfig) Build() (*RateTracker, error) {
	bucketDuration, err := getBucketDuration(c.TimeUnit)
	if err != nil {
		return nil, err
	}

	numBuckets, err := getNumBuckets(c.TimeWindow, bucketDuration)
	if err != nil {
		return nil, err
	}

	return NewRateTracker(numBuckets, bucketDuration), nil
}

// RateTracker is a struct that tracks a rate
type RateTracker struct {
	// buckets is a circular buffer of counts observed in each bucket
	buckets []int64

	// bucketDuration is the duration of each bucket
	bucketDuration time.Duration

	// cancel is the function to call to cancel the rate tracker
	cancel context.CancelFunc

	// wg is the wait group to wait for the rate tracker to stop
	wg sync.WaitGroup
}

// Add adds a count to the current bucket
func (r *RateTracker) Add(count int64) {
	r.buckets[0] += count
}

// GetRate returns the current rate
func (r *RateTracker) GetRate() float64 {
	var total int64
	for _, b := range r.buckets {
		total += b
	}
	return float64(total) / float64(len(r.buckets))
}

// ShiftBuckets shifts the buckets by one, dropping the oldest bucket and adding a new bucket
func (r *RateTracker) ShiftBuckets() {
	copy(r.buckets[1:], r.buckets)
	r.buckets[0] = 0
}

// Start starts the rate tracker
func (r *RateTracker) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	go r.HandleShift(ctx)
}

// HandleShift is a helper function to handle the shifting of buckets
func (r *RateTracker) HandleShift(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(r.bucketDuration)
	for {
		select {
		case <-ticker.C:
			r.ShiftBuckets()
		case <-ctx.Done():
			return
		}
	}
}

// Stop stops the rate tracker
func (r *RateTracker) Stop() {
	r.cancel()
	r.wg.Wait()
}

// NewRateTracker returns a new rate tracker
func NewRateTracker(bucketCount int, bucketDuration time.Duration) *RateTracker {
	return &RateTracker{
		buckets:        make([]int64, bucketCount),
		bucketDuration: bucketDuration,
	}
}

// getBucketDuration returns the duration of a bucket for a given time unit
func getBucketDuration(timeUnit TimeUnit) (time.Duration, error) {
	duration, ok := timeUnitMap[timeUnit]
	if !ok {
		return 0, fmt.Errorf("invalid time unit: %s", timeUnit)
	}
	return duration, nil
}

// getNumBuckets returns the number of buckets for a window of time and a bucket duration
func getNumBuckets(window time.Duration, bucketDuration time.Duration) (int, error) {
	if window < bucketDuration {
		return 0, errWindowTooSmall
	}
	return int(window / bucketDuration), nil
}

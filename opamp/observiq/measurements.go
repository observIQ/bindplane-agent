package observiq

import (
	"errors"
	"sync"
	"time"

	"github.com/golang/snappy"
	"github.com/observiq/bindplane-agent/internal/measurements"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type MeasurementsReporter interface {
	OTLPMeasurements() pmetric.Metrics
}

type NoopMeasurementsReporter struct{}

func (NoopMeasurementsReporter) OTLPMeasurements() pmetric.Metrics {
	return pmetric.NewMetrics()
}

// measurementsSender is a struct that handles periodically sending measurements via a custom message to an OpAMP endpoint.
type measurementsSender struct {
	logger      *zap.Logger
	reporter    MeasurementsReporter
	opampClient client.OpAMPClient
	interval    time.Duration

	changeIntervalChan chan time.Duration

	mux       *sync.Mutex
	isRunning bool
	done      chan struct{}
	wg        *sync.WaitGroup
}

func newMeasurementsSender(l *zap.Logger, reporter MeasurementsReporter, opampClient client.OpAMPClient, interval time.Duration) *measurementsSender {
	return &measurementsSender{
		logger:      l,
		reporter:    reporter,
		opampClient: opampClient,
		interval:    interval,

		changeIntervalChan: make(chan time.Duration, 1),
		mux:                &sync.Mutex{},
		isRunning:          false,
		done:               make(chan struct{}),
		wg:                 &sync.WaitGroup{},
	}
}

// Start starts the sender. It may be called multiple times, even if the sender is already started.
func (m *measurementsSender) Start() {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.isRunning {
		return
	}

	m.isRunning = true

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.loop()
	}()
}

// SetInterval changes the interval of the measurements sender.
func (m measurementsSender) SetInterval(d time.Duration) {
	select {
	case m.changeIntervalChan <- d:
	default:
		m.logger.Warn("Change interval chan was full, dropping change in reporting interval", zap.Duration("interval", d))
	}

}

func (m *measurementsSender) Stop() {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.isRunning {
		return
	}

	close(m.done)
	m.wg.Wait()

	m.isRunning = false
}

func (m *measurementsSender) loop() {
	t := newTicker()
	t.SetInterval(m.interval)
	defer t.Stop()

	for {
		select {
		case newInterval := <-m.changeIntervalChan:
			m.interval = newInterval
			t.SetInterval(newInterval)
		case <-m.done:
			return
		case <-t.Chan():
			if m.reporter == nil {
				// Continue if no reporter available
				continue
			}

			metrics := m.reporter.OTLPMeasurements()
			if metrics.DataPointCount() == 0 {
				// don't report empty payloads
				continue
			}

			// TODO: Make the same as bindplane metrics (share encoding and stuff)
			// Send metrics as snappy-encoded otlp proto
			marshaller := pmetric.ProtoMarshaler{}
			marshalled, err := marshaller.MarshalMetrics(metrics)
			if err != nil {
				m.logger.Error("Failed to marshal throughput metrics.", zap.Error(err))
				continue
			}

			encoded := snappy.Encode(nil, marshalled)

			cm := &protobufs.CustomMessage{
				Capability: measurements.ReportMeasurementsV1Capability,
				Type:       measurements.ReportMeasurementsType,
				Data:       encoded,
			}

			for {
				sendingChannel, err := m.opampClient.SendCustomMessage(cm)
				switch {
				case err == nil:
				case errors.Is(err, types.ErrCustomMessagePending):
					select {
					case <-sendingChannel:
						continue
					case <-m.done:
						return
					}
				default:
					m.logger.Error("Failed to report measurements", zap.Error(err))
				}
				break
			}
		}
	}
}

// ticker is essentially time.ticker, but it provides a SetInterval method
// that allows the interval to be changed. It also allows the interval
// to be configured to a negative or zero duration, in which case the ticker
// never fires.
type ticker struct {
	duration time.Duration
	ticker   *time.Ticker
}

func newTicker() *ticker {
	return &ticker{}
}

func (t *ticker) SetInterval(d time.Duration) {
	if t.duration == d {
		// Nothing to do, this is already the interval
		return
	}

	t.duration = d

	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}

	if d <= 0 {
		// Cannot make a ticker with zero or negative duration;
		// Attempts to use the channel will give a permanently blocking channel.
		return
	}

	t.ticker = time.NewTicker(d)
}

func (t *ticker) Chan() <-chan time.Time {
	if t.ticker == nil {
		// ticker never triggers if 0 or negative duration
		return make(<-chan time.Time)
	}
	return t.ticker.C
}

func (t *ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

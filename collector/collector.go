package collector

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/service"
)

// Collector wraps the open telemetry service.
type Collector struct {
	settings service.CollectorSettings
	running  bool
	cancel   context.CancelFunc
	mux      *sync.Mutex
	err      error
}

// New returns a new collector.
func New(settings service.CollectorSettings) *Collector {
	return &Collector{
		settings: settings,
		mux:      &sync.Mutex{},
	}
}

// Run will run the collector until the underlying service errors or completes.
func (c *Collector) Run(ctx context.Context) {
	// This mutex will ensure that only one caller can run the collector at a time.
	c.mux.Lock()
	defer c.mux.Unlock()

	// Toggle running state
	c.setRunning(true)
	defer c.setRunning(false)
	c.setError(nil)

	// Init a new service based on the collector settings.
	svc, err := service.New(c.settings)
	if err != nil {
		c.setError(err)
		return
	}

	// Use context to handle stopping the service from outside of this function.
	svcCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	c.cancel = cancel

	// Run the service in a separate go routine to avoid blocking context.
	runError := make(chan error)
	go func() {
		runError <- svc.Run()
	}()

	// Wait for context to cancel or a run error to occur.
	select {
	case <-svcCtx.Done():
		svc.Shutdown()
	case err := <-runError:
		c.setError(err)
	}
}

// Stop will stop the collector.
func (c *Collector) Stop() {
	if c.running && c.cancel != nil {
		c.cancel()
	}
}

// Running indicates if the collector is running.
func (c *Collector) Running() bool {
	return c.running
}

// setRunning will update the running state of the collector.
func (c *Collector) setRunning(state bool) {
	c.running = state
}

// Err returns the error state of the collector.
// TODO: Explore error channels instead to enable state listeners
func (c *Collector) Error() error {
	return c.err
}

// setError sets the error state of the collector.
func (c *Collector) setError(err error) {
	c.err = err
}

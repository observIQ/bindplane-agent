package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

// Collector is an interface for running the open telemetry collector.
type Collector interface {
	Run(context.Context) error
	Stop()
	Restart(context.Context) error
	Status() <-chan *Status
}

// collector is the standard implementation of the Collector interface.
type collector struct {
	configPath  string
	version     string
	loggingOpts []zap.Option
	mux         sync.Mutex
	svc         *service.Collector
	status      Status
	statusMux   sync.RWMutex
	statusChan  chan *Status
	startupChan chan error
	wg          *sync.WaitGroup
}

// New returns a new collector.
func New(configPath string, version string, loggingOpts []zap.Option) Collector {
	return &collector{
		configPath:  configPath,
		version:     version,
		loggingOpts: loggingOpts,
		statusChan:  make(chan *Status, 10),
		wg:          &sync.WaitGroup{},
	}
}

// Run will run the collector. This function will return an error
// if the collector was unable to startup.
func (c *collector) Run(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.svc != nil {
		return errors.New("service already running")
	}

	// The OT collector only supports using settings once during the lifetime
	// of a single collector insance. We must remake the settings on each startup.
	settings := NewSettings([]string{c.configPath}, c.version, c.loggingOpts)

	// The OT collector only supports calling run once during the lifetime
	// of a service. We must make a new instance each time we run the collector.
	svc, err := service.New(settings)
	if err != nil {
		err := fmt.Errorf("failed to create service: %w", err)
		c.sendStatus(false, err)
		return err
	}

	c.svc = svc
	c.startupChan = make(chan error, 1)
	c.wg = &sync.WaitGroup{}
	c.wg.Add(1)

	go c.runService(ctx)

	// A race condition exists in the OT collector where the shutdown channel
	// is not guaranteed to be initialized before the shutdown function is called.
	// We protect against this by waiting for startup to finish before unlocking the mutex.
	return c.waitForStartup(ctx)
}

// Stop will stop the collector.
func (c *collector) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.svc == nil {
		return
	}

	c.svc.Shutdown()
	c.wg.Wait()
	c.svc = nil
}

// Restart will restart the collector.
func (c *collector) Restart(ctx context.Context) error {
	c.Stop()
	return c.Run(ctx)
}

// runService will run the collector service. This is a blocking function
// that should be executed in a separate goroutine.
func (c *collector) runService(ctx context.Context) {
	defer c.wg.Done()
	err := c.svc.Run(ctx)
	c.sendStatus(false, err)

	if err != nil {
		c.startupChan <- err
	}
}

// waitForStartup waits for the service to startup before exiting.
func (c *collector) waitForStartup(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()

	for {
		if c.svc.GetState() == service.Running {
			c.sendStatus(true, nil)
			return nil
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		case err := <-c.startupChan:
			return err
		}
	}
}

// Status will return the status of the collector.
func (c *collector) Status() <-chan *Status {
	return c.statusChan
}

// sendStatus will set the status of the collector
func (c *collector) sendStatus(running bool, err error) {
	select {
	case c.statusChan <- &Status{running, err}:
	default:
	}
}

// Status is the status of a collector.
type Status struct {
	Running bool
	Err     error
}

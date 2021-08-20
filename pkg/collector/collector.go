package collector

import (
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

// Collector wraps the open telemetry collector.
type Collector struct {
	configPath  string
	settings    service.CollectorSettings
	svc         *service.Collector
	svcMux      *sync.Mutex
	statusChan  chan *Status
	startupChan chan error
	wg          *sync.WaitGroup
}

// New returns a new collector.
func New(configPath string, version string, loggingOpts []zap.Option) *Collector {
	settings := NewSettings(configPath, version, loggingOpts)
	return &Collector{
		configPath: configPath,
		settings:   settings,
		svcMux:     &sync.Mutex{},
		statusChan: make(chan *Status, 10),
		wg:         &sync.WaitGroup{},
	}
}

// Run will run the collector. This function will return an error
// if the collector was unable to startup.
func (c *Collector) Run() error {
	c.svcMux.Lock()
	defer c.svcMux.Unlock()

	if c.svc != nil {
		return errors.New("service already running")
	}

	// The OT collector only supports calling run once during the lifetime
	// of a service. We must make a new instance each time we run the collector.
	svc, err := service.New(c.settings)
	if err != nil {
		c.sendStatus(false, err)
		return fmt.Errorf("failed to init service: %w", err)
	}

	c.svc = svc
	c.startupChan = make(chan error, 1)
	c.wg = &sync.WaitGroup{}
	c.wg.Add(1)

	go c.runCollector()

	// A race condition exists in the OT collector where the shutdown channel
	// is not guaranteed to be initialized before the shutdown function is called.
	// We protect against this by waiting for startup to finish before unlocking the mutex.
	return c.waitForStartup()
}

// Stop will stop the collector.
func (c *Collector) Stop() {
	c.svcMux.Lock()
	defer c.svcMux.Unlock()

	if c.svc == nil {
		return
	}

	c.svc.Shutdown()
	c.wg.Wait()
	c.svc = nil
}

// Restart will restart the collector.
func (c *Collector) Restart() error {
	c.Stop()
	return c.Run()
}

// runCollector will run the collector. This is a blocking function
// that should be executed in a separate goroutine.
func (c *Collector) runCollector() {
	defer c.wg.Done()

	err := c.svc.Run()
	c.sendStatus(false, err)

	if err != nil {
		c.startupChan <- err
	}
}

// waitForStartup waits for the service to startup before exiting.
func (c *Collector) waitForStartup() error {
	for {
		select {
		// If the service is able to transmit a running state, this means
		// that initial service startup was successful.
		case state := <-c.svc.GetStateChannel():
			if state == service.Running {
				c.sendStatus(true, nil)
				return nil
			}
		// A flaw exists in the OT startup function where an early error may occur
		// without sending any states through the state channel. To handle this,
		// we use an errChan to exit immediately.
		case err := <-c.startupChan:
			return err
		}
	}
}

// ConfigPath will return the config path of the collector.
func (c *Collector) ConfigPath() string {
	return c.configPath
}

// ValidateConfig will validate the collector's config.
func (c *Collector) ValidateConfig() error {
	_, err := configtest.LoadConfigAndValidate(c.configPath, c.settings.Factories)
	return err
}

// Status will return the status of the collector.
func (c *Collector) Status() <-chan *Status {
	return c.statusChan
}

// sendStatus will send a status through the status channel.
// If the channel is full, the status is discarded.
func (c *Collector) sendStatus(running bool, err error) {
	status := &Status{Running: running, Err: err}

	select {
	case c.statusChan <- status:
	default:
	}
}

// Status is the status of a collector.
type Status struct {
	Running bool
	Err     error
}

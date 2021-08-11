package context

import (
	"context"
	"os"
	"os/signal"
)

// WithInterrupt returns a context that cancels when an interrupt signal is received.
func WithInterrupt() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	cancelFunc := func() {
		signal.Stop(signalChan)
		cancel()
	}

	return ctx, cancelFunc
}

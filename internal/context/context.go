package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/observiq/observiq-collector/internal/process"
)

// EmptyContext returns an empty context
func EmptyContext() context.Context {
	return context.Background()
}

// WithParent returns a context that cancels when the supplied parent process exits.
func WithParent(ppid int) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if !process.MatchesParent(ppid) {
					cancel()
				}
			}
		}
	}()

	return ctx
}

// WithInterrupt returns a context that cancels when an interrupt signal is received.
func WithInterrupt(ctx context.Context) (context.Context, context.CancelFunc) {
	interruptCtx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-interruptCtx.Done():
		}
	}()

	cancelFunc := func() {
		signal.Stop(signalChan)
		cancel()
	}

	return interruptCtx, cancelFunc
}

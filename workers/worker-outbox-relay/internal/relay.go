package internal

import (
	"context"
	"log/slog"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

type outboxFlusher interface {
	Flush(ctx context.Context, bus queue.Bus) error
}

type Relay struct {
	outbox   outboxFlusher
	bus      queue.Bus
	interval time.Duration
	logger   *slog.Logger
}

func NewRelay(outbox outboxFlusher, bus queue.Bus, interval time.Duration, logger *slog.Logger) *Relay {
	if interval <= 0 {
		interval = 1 * time.Second
	}
	return &Relay{
		outbox:   outbox,
		bus:      bus,
		interval: interval,
		logger:   logger,
	}
}

func (r *Relay) FlushOnce(ctx context.Context) error {
	return r.outbox.Flush(ctx, r.bus)
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		if err := r.FlushOnce(ctx); err != nil {
			r.logger.Error("outbox flush failed", "worker", "worker-outbox-relay", "error", err.Error())
		}
		select {
		case <-ctx.Done():
			r.logger.Info("worker stopped", "worker", "worker-outbox-relay")
			return
		case <-ticker.C:
		}
	}
}

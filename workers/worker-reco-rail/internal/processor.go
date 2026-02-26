package internal

import (
	"context"
	"log/slog"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

type Processor struct {
	name          string
	expectedTopic string
	guard         *queue.IdempotencyGuard
	logger        *slog.Logger
}

func NewProcessor(name, expectedTopic string, logger *slog.Logger) *Processor {
	return &Processor{name: name, expectedTopic: expectedTopic, guard: queue.NewScopedIdempotencyGuard(name), logger: logger}
}

func (p *Processor) Handle(_ context.Context, event queue.Event) error {
	if p.guard.Seen(event.ID) {
		p.logger.Info("duplicate event ignored", "worker", p.name, "event_id", event.ID)
		return nil
	}
	if event.Topic != p.expectedTopic {
		p.logger.Info("unexpected topic skipped", "worker", p.name, "topic", event.Topic)
		return nil
	}
	p.logger.Info("event processed", "worker", p.name, "topic", event.Topic, "bytes", len(event.Payload))
	return nil
}

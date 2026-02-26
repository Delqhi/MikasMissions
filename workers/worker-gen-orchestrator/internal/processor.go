package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/generatorrunstore"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

type Processor struct {
	bus    queue.Bus
	guard  *queue.IdempotencyGuard
	logger *slog.Logger
}

func NewProcessor(bus queue.Bus, logger *slog.Logger) *Processor {
	return &Processor{bus: bus, guard: queue.NewScopedIdempotencyGuard("worker-gen-orchestrator"), logger: logger}
}

func (p *Processor) Topic() string {
	return "video.run.requested.v1"
}

func (p *Processor) Consumer() string {
	return "worker-gen-orchestrator"
}

func (p *Processor) Handle(ctx context.Context, event queue.Event) error {
	if p.guard.Seen(event.ID) {
		p.logger.Info("duplicate event ignored", "worker", p.Consumer(), "event_id", event.ID)
		return nil
	}
	if event.Topic != p.Topic() {
		p.logger.Info("unexpected topic skipped", "worker", p.Consumer(), "topic", event.Topic)
		return nil
	}
	var incoming contractsevents.VideoRunRequestedV1
	if err := json.Unmarshal(event.Payload, &incoming); err != nil {
		return err
	}
	if err := incoming.Validate(); err != nil {
		return err
	}
	payload, err := json.Marshal(contractsevents.VideoRunStepCompletedV1{
		RunID:       incoming.RunID,
		Step:        "orchestrator",
		Status:      "completed",
		Details:     "generation pipeline scheduled",
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	if err := p.bus.Publish(ctx, queue.Event{
		ID:      uuid.NewString(),
		Topic:   "video.run.step.completed.v1",
		Payload: payload,
	}); err != nil {
		return err
	}
	_ = generatorrunstore.SetRunStatus(ctx, incoming.RunID, "running", "")
	_ = generatorrunstore.AppendRunLog(ctx, incoming.RunID, "orchestrator", "completed", "generation pipeline scheduled")
	p.logger.Info("run orchestrated", "worker", p.Consumer(), "run_id", incoming.RunID)
	return nil
}

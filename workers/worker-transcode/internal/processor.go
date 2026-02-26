package internal

import (
	"context"
	"encoding/json"
	"log/slog"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/pipeline"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

type Processor struct {
	bus    queue.Bus
	guard  *queue.IdempotencyGuard
	logger *slog.Logger
}

func NewProcessor(bus queue.Bus, logger *slog.Logger) *Processor {
	return &Processor{bus: bus, guard: queue.NewScopedIdempotencyGuard("worker-transcode"), logger: logger}
}

func (p *Processor) Topic() string {
	return "media.transcode.requested.v1"
}

func (p *Processor) Consumer() string {
	return "worker-transcode"
}

func (p *Processor) Handle(ctx context.Context, event queue.Event) error {
	if p.guard.Seen(event.ID) {
		p.logger.Info("duplicate event ignored", "worker", "worker-transcode", "event_id", event.ID)
		return nil
	}
	if event.Topic != p.Topic() {
		p.logger.Info("unexpected topic skipped", "worker", "worker-transcode", "topic", event.Topic)
		return nil
	}
	var incoming contractsevents.MediaTranscodeRequestedV1
	if err := json.Unmarshal(event.Payload, &incoming); err != nil {
		return err
	}
	if err := incoming.Validate(); err != nil {
		return err
	}
	outgoing, err := pipeline.BuildTranscodedMedia(incoming)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(outgoing)
	if err != nil {
		return err
	}
	if err := p.bus.Publish(ctx, queue.Event{
		ID:      uuid.NewString(),
		Topic:   "media.transcoded.v1",
		Payload: payload,
	}); err != nil {
		return err
	}
	p.logger.Info("event processed", "worker", "worker-transcode", "asset_id", incoming.AssetID)
	return nil
}

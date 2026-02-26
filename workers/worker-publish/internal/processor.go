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

type catalogProjector interface {
	ProjectEpisode(ctx context.Context, req episodeProjectionRequest) error
}

type Processor struct {
	bus       queue.Bus
	guard     *queue.IdempotencyGuard
	projector catalogProjector
	logger    *slog.Logger
}

func NewProcessor(bus queue.Bus, logger *slog.Logger) *Processor {
	return &Processor{
		bus:       bus,
		guard:     queue.NewScopedIdempotencyGuard("worker-publish"),
		projector: NewCatalogProjectorFromEnv(),
		logger:    logger,
	}
}

func (p *Processor) Topic() string {
	return "media.approved.v1"
}

func (p *Processor) Consumer() string {
	return "worker-publish"
}

func (p *Processor) Handle(ctx context.Context, event queue.Event) error {
	if p.guard.Seen(event.ID) {
		p.logger.Info("duplicate event ignored", "worker", "worker-publish", "event_id", event.ID)
		return nil
	}
	if event.Topic != p.Topic() {
		p.logger.Info("unexpected topic skipped", "worker", "worker-publish", "topic", event.Topic)
		return nil
	}
	var incoming contractsevents.MediaApprovedV1
	if err := json.Unmarshal(event.Payload, &incoming); err != nil {
		return err
	}
	if err := incoming.Validate(); err != nil {
		return err
	}
	outgoing, err := pipeline.BuildEpisodePublished(incoming)
	if err != nil {
		return err
	}
	if err := p.projector.ProjectEpisode(ctx, episodeProjectionRequest{
		EpisodeID:     outgoing.EpisodeID,
		AgeBand:       outgoing.AgeBand,
		LearningTags:  outgoing.LearningTags,
		PlaybackReady: true,
	}); err != nil {
		return err
	}
	payload, err := json.Marshal(outgoing)
	if err != nil {
		return err
	}
	if err := p.bus.Publish(ctx, queue.Event{
		ID:      uuid.NewString(),
		Topic:   "episode.published.v1",
		Payload: payload,
	}); err != nil {
		return err
	}
	p.logger.Info("event processed", "worker", "worker-publish", "asset_id", incoming.AssetID, "episode_id", outgoing.EpisodeID)
	return nil
}

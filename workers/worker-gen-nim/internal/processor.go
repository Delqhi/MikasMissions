package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/generatorprovider"
	"github.com/delqhi/mikasmissions/platform/libs/generatorrunstore"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

type Processor struct {
	bus          queue.Bus
	guard        *queue.IdempotencyGuard
	profileStore modelProfileReader
	logger       *slog.Logger
}

func NewProcessor(bus queue.Bus, logger *slog.Logger) *Processor {
	return &Processor{
		bus:          bus,
		guard:        queue.NewScopedIdempotencyGuard("worker-gen-nim"),
		profileStore: newModelProfileReaderFromEnv(),
		logger:       logger,
	}
}

func (p *Processor) Topic() string {
	return "video.run.requested.v1"
}

func (p *Processor) Consumer() string {
	return "worker-gen-nim"
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
	profile, err := p.profileStore.GetProfile(incoming.ModelProfileID)
	if err != nil {
		p.publishFailed(ctx, incoming.RunID, "nim_profile_error", err.Error())
		return nil
	}
	provider, err := generatorprovider.NewProvider(profile)
	if err != nil {
		p.publishFailed(ctx, incoming.RunID, "nim_provider_error", err.Error())
		return nil
	}
	result, err := provider.GenerateVideo(generatorprovider.GenerateRequest{
		RunID:        incoming.RunID,
		InputPayload: incoming.InputPayload,
	})
	if err != nil {
		p.publishFailed(ctx, incoming.RunID, "nim_provider_error", err.Error())
		return nil
	}
	ready := contractsevents.VideoAssetReadyV1{
		RunID:              incoming.RunID,
		AssetID:            result.AssetID,
		SourceURL:          result.SourceURL,
		DurationMS:         result.DurationMS,
		ContentSuitability: incoming.ContentSuitability,
		AgeBand:            incoming.AgeBand,
		UploaderID:         incoming.RequestedBy,
		ReadyAt:            time.Now().UTC().Format(time.RFC3339),
	}
	readyPayload, err := json.Marshal(ready)
	if err != nil {
		return err
	}
	if err := p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "video.asset.ready.v1", Payload: readyPayload}); err != nil {
		return err
	}
	stepPayload, err := json.Marshal(contractsevents.VideoRunStepCompletedV1{
		RunID:       incoming.RunID,
		Step:        "nim",
		Status:      "completed",
		Details:     "nim generation completed",
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	if err := p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "video.run.step.completed.v1", Payload: stepPayload}); err != nil {
		return err
	}
	_ = generatorrunstore.AppendRunLog(ctx, incoming.RunID, "nim", "completed", "nim generation completed")
	p.logger.Info("nim generation completed", "run_id", incoming.RunID, "asset_id", result.AssetID)
	return nil
}

func (p *Processor) publishFailed(ctx context.Context, runID, code, message string) {
	_ = generatorrunstore.SetRunStatus(ctx, runID, "failed", message)
	_ = generatorrunstore.AppendRunLog(ctx, runID, "nim", "failed", message)
	payload, err := json.Marshal(contractsevents.VideoRunFailedV1{
		RunID:        runID,
		Step:         "nim",
		ErrorCode:    code,
		ErrorMessage: message,
		FailedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		p.logger.Error("failed to encode failed event", "error", err.Error(), "run_id", runID)
		return
	}
	if err := p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "video.run.failed.v1", Payload: payload}); err != nil {
		p.logger.Error("failed to publish failed event", "error", err.Error(), "run_id", runID)
	}
}

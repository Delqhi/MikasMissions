package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
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
	return &Processor{bus: bus, guard: queue.NewScopedIdempotencyGuard("worker-gen-qc"), logger: logger}
}

func (p *Processor) Topic() string {
	return "video.asset.ready.v1"
}

func (p *Processor) Consumer() string {
	return "worker-gen-qc"
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
	var incoming contractsevents.VideoAssetReadyV1
	if err := json.Unmarshal(event.Payload, &incoming); err != nil {
		return err
	}
	if err := incoming.Validate(); err != nil {
		return err
	}
	if err := qcValidate(incoming); err != nil {
		return p.publishFailed(ctx, incoming, "generation_qc_failed", err.Error())
	}
	if err := p.publishMediaUploaded(ctx, incoming); err != nil {
		return err
	}
	if err := p.publishStepCompleted(ctx, incoming.RunID); err != nil {
		return err
	}
	_ = generatorrunstore.SetRunStatus(ctx, incoming.RunID, "publish_queued", "")
	_ = generatorrunstore.AppendRunLog(ctx, incoming.RunID, "qc", "completed", "qc checks passed and upload queued")
	p.logger.Info("qc passed and media upload event emitted", "run_id", incoming.RunID, "asset_id", incoming.AssetID)
	return nil
}

func qcValidate(event contractsevents.VideoAssetReadyV1) error {
	if !strings.HasPrefix(event.SourceURL, "http://") && !strings.HasPrefix(event.SourceURL, "https://") {
		return errInvalidSourceURL
	}
	if event.DurationMS < 30000 || event.DurationMS > 1800000 {
		return errInvalidDuration
	}
	return nil
}

func (p *Processor) publishMediaUploaded(ctx context.Context, incoming contractsevents.VideoAssetReadyV1) error {
	payload, err := json.Marshal(contractsevents.MediaUploadedV1{
		AssetID:   incoming.AssetID,
		SourceURL: incoming.SourceURL,
		Uploader:  uploaderFromEvent(incoming),
		TraceID:   incoming.RunID,
	})
	if err != nil {
		return err
	}
	return p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "media.uploaded.v1", Payload: payload})
}

func uploaderFromEvent(event contractsevents.VideoAssetReadyV1) string {
	if event.UploaderID != "" {
		return event.UploaderID
	}
	return "admin-studio"
}

func (p *Processor) publishStepCompleted(ctx context.Context, runID string) error {
	payload, err := json.Marshal(contractsevents.VideoRunStepCompletedV1{
		RunID:       runID,
		Step:        "qc",
		Status:      "completed",
		Details:     "qc checks passed and upload queued",
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "video.run.step.completed.v1", Payload: payload})
}

func (p *Processor) publishFailed(ctx context.Context, incoming contractsevents.VideoAssetReadyV1, code, message string) error {
	_ = generatorrunstore.SetRunStatus(ctx, incoming.RunID, "failed", message)
	_ = generatorrunstore.AppendRunLog(ctx, incoming.RunID, "qc", "failed", message)
	payload, err := json.Marshal(contractsevents.VideoRunFailedV1{
		RunID:        incoming.RunID,
		Step:         "qc",
		ErrorCode:    code,
		ErrorMessage: message,
		FailedAt:     time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return p.bus.Publish(ctx, queue.Event{ID: uuid.NewString(), Topic: "video.run.failed.v1", Payload: payload})
}

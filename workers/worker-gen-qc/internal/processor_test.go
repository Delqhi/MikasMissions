package internal

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func TestProcessorPublishesMediaUploadedOnQCPass(t *testing.T) {
	bus := queue.NewInMemoryBus()
	processor := NewProcessor(bus, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	uploaded := 0
	if err := bus.Subscribe(context.Background(), "media.uploaded.v1", "test-uploaded", func(_ context.Context, event queue.Event) error {
		uploaded++
		return nil
	}); err != nil {
		t.Fatalf("subscribe uploaded: %v", err)
	}

	payload, err := json.Marshal(contractsevents.VideoAssetReadyV1{
		RunID:              "run-1",
		AssetID:            "asset-1",
		SourceURL:          "https://cdn.example/asset-1.mp4",
		DurationMS:         120000,
		ContentSuitability: "core",
		AgeBand:            "6-11",
		UploaderID:         "admin-1",
		ReadyAt:            "2026-03-11T10:00:00Z",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := processor.Handle(context.Background(), queue.Event{ID: "e1", Topic: processor.Topic(), Payload: payload}); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if uploaded != 1 {
		t.Fatalf("expected 1 uploaded event, got %d", uploaded)
	}
}

func TestProcessorPublishesFailureOnQCReject(t *testing.T) {
	bus := queue.NewInMemoryBus()
	processor := NewProcessor(bus, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	failed := 0
	if err := bus.Subscribe(context.Background(), "video.run.failed.v1", "test-failed", func(_ context.Context, event queue.Event) error {
		failed++
		return nil
	}); err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	payload, err := json.Marshal(contractsevents.VideoAssetReadyV1{
		RunID:              "run-1",
		AssetID:            "asset-1",
		SourceURL:          "ftp://invalid",
		DurationMS:         120000,
		ContentSuitability: "core",
		AgeBand:            "6-11",
		UploaderID:         "admin-1",
		ReadyAt:            "2026-03-11T10:00:00Z",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := processor.Handle(context.Background(), queue.Event{ID: "e2", Topic: processor.Topic(), Payload: payload}); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if failed != 1 {
		t.Fatalf("expected 1 failed event, got %d", failed)
	}
}

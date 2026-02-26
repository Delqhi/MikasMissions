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

func TestProcessorIdempotency(t *testing.T) {
	t.Setenv("CATALOG_URL", "")
	bus := queue.NewInMemoryBus()
	processor := NewProcessor(bus, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	var seen int
	_ = bus.Subscribe(context.Background(), "episode.published.v1", "test-publish", func(_ context.Context, event queue.Event) error {
		var decoded contractsevents.EpisodePublishedV1
		_ = json.Unmarshal(event.Payload, &decoded)
		seen++
		return nil
	})
	payload, _ := json.Marshal(contractsevents.MediaApprovedV1{
		AssetID:      "asset-1",
		AgeBand:      "6-11",
		LearningTags: []string{"farben"},
	})
	event := queue.Event{ID: "evt1", Topic: "media.approved.v1", Payload: payload}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("first handle failed: %v", err)
	}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("duplicate handle failed: %v", err)
	}
	if seen != 1 {
		t.Fatalf("expected 1 published event, got %d", seen)
	}
}

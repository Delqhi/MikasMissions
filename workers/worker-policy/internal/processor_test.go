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
	bus := queue.NewInMemoryBus()
	processor := NewProcessor(bus, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	var reviewedSeen int
	var approvedSeen int
	_ = bus.Subscribe(context.Background(), "media.reviewed.v1", "test-reviewed", func(_ context.Context, event queue.Event) error {
		var decoded contractsevents.MediaReviewedV1
		_ = json.Unmarshal(event.Payload, &decoded)
		reviewedSeen++
		return nil
	})
	_ = bus.Subscribe(context.Background(), "media.approved.v1", "test-approved", func(_ context.Context, event queue.Event) error {
		var decoded contractsevents.MediaApprovedV1
		_ = json.Unmarshal(event.Payload, &decoded)
		approvedSeen++
		return nil
	})
	payload, _ := json.Marshal(contractsevents.MediaTranscodedV1{
		AssetID: "asset-1",
		Renditions: []contractsevents.MediaRendition{
			{Profile: "720p", URL: "https://cdn.local/asset-1/720p.m3u8"},
		},
		DurationMS: 120000,
	})
	event := queue.Event{ID: "evt1", Topic: "media.transcoded.v1", Payload: payload}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("first handle failed: %v", err)
	}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("duplicate handle failed: %v", err)
	}
	if reviewedSeen != 1 || approvedSeen != 1 {
		t.Fatalf("expected reviewed=1 and approved=1, got %d and %d", reviewedSeen, approvedSeen)
	}
}

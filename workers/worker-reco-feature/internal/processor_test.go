package internal

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func TestProcessorIdempotency(t *testing.T) {
	processor := NewProcessor("worker-reco-feature", "watch.event.v1", slog.New(slog.NewJSONHandler(io.Discard, nil)))
	event := queue.Event{ID: "evt1", Topic: "watch.event.v1", Payload: []byte("{}")}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("first handle failed: %v", err)
	}
	if err := processor.Handle(context.Background(), event); err != nil {
		t.Fatalf("duplicate handle failed: %v", err)
	}
}

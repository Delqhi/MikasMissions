package internal

import (
	"context"
	"encoding/json"
	"testing"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func TestUploadAssetPublishesMediaUploadedEvent(t *testing.T) {
	bus := queue.NewInMemoryBus()
	service, err := NewService(bus)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	defer func() {
		_ = service.Close()
	}()
	var seen int
	if err := bus.Subscribe(context.Background(), "media.uploaded.v1", "test-creator", func(_ context.Context, event queue.Event) error {
		var decoded contractsevents.MediaUploadedV1
		if err := json.Unmarshal(event.Payload, &decoded); err != nil {
			t.Fatalf("unmarshal event: %v", err)
		}
		seen++
		return nil
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	if _, err := service.UploadAsset(context.Background(), UploadRequest{SourceURL: "https://cdn.local/a.mp4", UploaderID: "u-1"}); err != nil {
		t.Fatalf("upload asset: %v", err)
	}
	if seen != 1 {
		t.Fatalf("expected 1 event, got %d", seen)
	}
}

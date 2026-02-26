package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestWatchEventV1Contract(t *testing.T) {
	raw := []byte(`{"child_profile_id":"c1","episode_id":"ep1","watch_ms":1000,"event_time":"2026-03-01T12:00:00Z"}`)
	var event WatchEventV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

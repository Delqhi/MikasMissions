package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestChildSessionCappedV1Contract(t *testing.T) {
	raw := []byte(`{"child_profile_id":"c1","session_id":"s1","capped_at":"2026-03-01T10:00:00Z","reason":"limit_reached","limit_minutes":30}`)
	var event ChildSessionCappedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

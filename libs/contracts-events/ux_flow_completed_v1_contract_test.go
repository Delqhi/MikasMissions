package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestUXFlowCompletedV1Contract(t *testing.T) {
	raw := []byte(`{"flow_id":"f1","parent_user_id":"p1","child_profile_id":"c1","step":"signup","status":"done","completed_at":"2026-03-01T10:00:00Z"}`)
	var event UXFlowCompletedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

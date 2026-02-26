package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestParentGateChallengeV1Contract(t *testing.T) {
	raw := []byte(`{"gate_id":"g1","parent_user_id":"p1","child_profile_id":"c1","method":"pin","verified":true,"challenge_time":"2026-03-01T10:00:00Z"}`)
	var event ParentGateChallengeV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

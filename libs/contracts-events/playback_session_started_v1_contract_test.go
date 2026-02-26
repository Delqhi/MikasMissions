package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestPlaybackSessionStartedV1Contract(t *testing.T) {
	raw := []byte(`{"playback_session_id":"s1","child_profile_id":"c1","episode_id":"ep1","started_at":"2026-03-01T10:00:00Z","safety_mode":"strict"}`)
	var event PlaybackSessionStartedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestPlaybackSessionStartedV1ContractRejectsMissingSafety(t *testing.T) {
	event := PlaybackSessionStartedV1{
		PlaybackSessionID: "s1",
		ChildProfileID:    "c1",
		EpisodeID:         "ep1",
		StartedAt:         "2026-03-01T10:00:00Z",
	}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

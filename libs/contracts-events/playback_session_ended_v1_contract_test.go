package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestPlaybackSessionEndedV1Contract(t *testing.T) {
	raw := []byte(`{"playback_session_id":"s1","child_profile_id":"c1","episode_id":"ep1","ended_at":"2026-03-01T10:15:00Z","watched_ms":120000,"capped":false}`)
	var event PlaybackSessionEndedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestPlaybackSessionEndedV1ContractRejectsNegativeWatch(t *testing.T) {
	event := PlaybackSessionEndedV1{
		PlaybackSessionID: "s1",
		ChildProfileID:    "c1",
		EpisodeID:         "ep1",
		EndedAt:           "2026-03-01T10:15:00Z",
		WatchedMS:         -1,
	}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

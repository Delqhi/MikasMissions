package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestEpisodePublishedV1Contract(t *testing.T) {
	raw := []byte(`{"episode_id":"ep1","age_band":"3-5","learning_tags":["colors","letters"]}`)
	var event EpisodePublishedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

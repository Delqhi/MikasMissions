package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestSafetyFilterAppliedV1Contract(t *testing.T) {
	raw := []byte(`{"child_profile_id":"c1","episode_id":"ep1","safety_mode":"strict","filters":["no_chat","no_external_links"],"reason":"age_policy"}`)
	var event SafetyFilterAppliedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

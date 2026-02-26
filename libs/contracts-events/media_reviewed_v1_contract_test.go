package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestMediaReviewedV1Contract(t *testing.T) {
	raw := []byte(`{"asset_id":"a1","policy_result":"approved","age_band":"6-11"}`)
	var event MediaReviewedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

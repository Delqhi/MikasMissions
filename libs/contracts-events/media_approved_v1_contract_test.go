package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestMediaApprovedV1Contract(t *testing.T) {
	raw := []byte(`{"asset_id":"a1","age_band":"6-11","learning_tags":["farben"]}`)
	var event MediaApprovedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

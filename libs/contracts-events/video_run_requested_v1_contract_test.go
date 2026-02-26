package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestVideoRunRequestedV1Contract(t *testing.T) {
	raw := []byte(`{"run_id":"run1","workflow_id":"wf1","model_profile_id":"nim-default","input_payload":{"theme":"space"},"auto_publish":false,"priority":"normal","content_suitability":"core","age_band":"6-11","requested_by":"admin1","requested_at":"2026-03-11T09:00:00Z","trace_id":"run1"}`)
	var event VideoRunRequestedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestVideoRunRequestedV1ContractRejectsMissingMetadata(t *testing.T) {
	event := VideoRunRequestedV1{RunID: "run1", WorkflowID: "wf1", ModelProfileID: "nim-default", Priority: "normal", RequestedBy: "admin1"}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

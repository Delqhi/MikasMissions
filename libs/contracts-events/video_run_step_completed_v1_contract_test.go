package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestVideoRunStepCompletedV1Contract(t *testing.T) {
	raw := []byte(`{"run_id":"run1","step":"nim","status":"completed","details":"asset generated","completed_at":"2026-03-11T09:00:00Z"}`)
	var event VideoRunStepCompletedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestVideoRunStepCompletedV1ContractRejectsMissingStep(t *testing.T) {
	event := VideoRunStepCompletedV1{RunID: "run1", Status: "completed", CompletedAt: "2026-03-11T09:00:00Z"}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

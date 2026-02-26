package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestVideoRunFailedV1Contract(t *testing.T) {
	raw := []byte(`{"run_id":"run1","step":"nim","error_code":"nim_provider_error","error_message":"timeout","failed_at":"2026-03-11T09:00:00Z"}`)
	var event VideoRunFailedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestVideoRunFailedV1ContractRejectsMissingError(t *testing.T) {
	event := VideoRunFailedV1{RunID: "run1", Step: "nim", FailedAt: "2026-03-11T09:00:00Z"}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

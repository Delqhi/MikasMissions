package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestVideoWorkflowCreatedV1Contract(t *testing.T) {
	raw := []byte(`{"workflow_id":"wf1","version":1,"created_by":"admin1","created_at":"2026-03-11T09:00:00Z"}`)
	var event VideoWorkflowCreatedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestVideoWorkflowCreatedV1ContractRejectsInvalidVersion(t *testing.T) {
	event := VideoWorkflowCreatedV1{WorkflowID: "wf1", CreatedBy: "admin1", CreatedAt: "2026-03-11T09:00:00Z", Version: 0}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestParentControlsUpdatedV1Contract(t *testing.T) {
	raw := []byte(`{"parent_user_id":"p1","child_profile_id":"c1","safety_mode":"strict","session_limit_minutes":45,"external_links":false,"audit_event_id":"audit-1"}`)
	var event ParentControlsUpdatedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestParentControlsUpdatedV1ContractRejectsInvalidLimit(t *testing.T) {
	event := ParentControlsUpdatedV1{
		ParentUserID:     "p1",
		ChildProfileID:   "c1",
		SafetyMode:       "strict",
		SessionLimitMins: 0,
		AuditEventID:     "audit-1",
	}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

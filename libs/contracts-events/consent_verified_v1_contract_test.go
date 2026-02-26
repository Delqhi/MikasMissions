package contractsevents

import (
	"encoding/json"
	"testing"
)

func TestConsentVerifiedV1Contract(t *testing.T) {
	raw := []byte(`{"consent_id":"c1","parent_user_id":"p1","method":"card","verified_at":"2026-03-01T10:00:00Z"}`)
	var event ConsentVerifiedV1
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestConsentVerifiedV1ContractRejectsMissingField(t *testing.T) {
	event := ConsentVerifiedV1{
		ConsentID:    "c1",
		ParentUserID: "p1",
		Method:       "card",
	}
	if err := event.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

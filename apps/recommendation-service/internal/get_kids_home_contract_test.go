package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestGetKidsHomeContractHasCanonicalFields(t *testing.T) {
	service := NewService()
	handler := GetKidsHome(service)

	req := httptest.NewRequest(http.MethodGet, "/v1/kids/home?child_profile_id=child-1&mode=early", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp contractsapi.KidsHomeResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.ChildProfileID == "" {
		t.Fatalf("expected child_profile_id")
	}
	if resp.Mode != "early" {
		t.Fatalf("expected early mode, got %q", resp.Mode)
	}
	if resp.SafetyMode != contractsapi.SafetyModeStrict {
		t.Fatalf("expected strict safety mode, got %q", resp.SafetyMode)
	}
	if len(resp.PrimaryActions) == 0 {
		t.Fatalf("expected primary actions")
	}
	if len(resp.Rails) == 0 {
		t.Fatalf("expected rails")
	}
	first := resp.Rails[0]
	if first.EpisodeID == "" || first.Summary == "" || first.ThumbnailURL == "" {
		t.Fatalf("expected canonical rail fields")
	}
	if first.DurationMS <= 0 {
		t.Fatalf("expected positive duration_ms")
	}
	if first.ReasonCode == "" {
		t.Fatalf("expected reason_code")
	}
}

func TestGetKidsHomeContractRejectsMissingChildProfile(t *testing.T) {
	service := NewService()
	handler := GetKidsHome(service)

	req := httptest.NewRequest(http.MethodGet, "/v1/kids/home?mode=core", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	var apiErr contractsapi.APIError
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("unmarshal api error: %v", err)
	}
	if apiErr.Code != "missing_child_profile" {
		t.Fatalf("expected missing_child_profile, got %q", apiErr.Code)
	}
}

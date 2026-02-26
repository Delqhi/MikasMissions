package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestPutParentControlsContractRejectsInvalidSessionLimit(t *testing.T) {
	store := NewStore()
	mux := NewMux(store)

	req := httptest.NewRequest(
		http.MethodPut,
		"/v1/parents/controls/child-1?parent_user_id=parent-1",
		strings.NewReader(`{"autoplay":false,"chat_enabled":false,"external_links":false,"session_limit_minutes":1,"bedtime_window":"20:00-07:00","safety_mode":"strict"}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	var apiErr contractsapi.APIError
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("unmarshal api error: %v", err)
	}
	if apiErr.Code != "invalid_session_limit" {
		t.Fatalf("expected invalid_session_limit, got %q", apiErr.Code)
	}
}

func TestPutParentControlsContractReturnsControlsEnvelope(t *testing.T) {
	store := NewStore()
	mux := NewMux(store)

	req := httptest.NewRequest(
		http.MethodPut,
		"/v1/parents/controls/child-1?parent_user_id=parent-1",
		strings.NewReader(`{"autoplay":false,"chat_enabled":false,"external_links":false,"session_limit_minutes":40,"bedtime_window":"20:00-07:00","safety_mode":"strict"}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp contractsapi.ParentControlsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.ChildProfileID != "child-1" {
		t.Fatalf("expected child-1, got %q", resp.ChildProfileID)
	}
	if resp.Controls.SessionLimitMinutes != 40 {
		t.Fatalf("expected session limit 40, got %d", resp.Controls.SessionLimitMinutes)
	}
}

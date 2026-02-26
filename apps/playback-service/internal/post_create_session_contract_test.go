package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestPostCreateSessionContractSuccess(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/playback/sessions",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"tablet","entitlement_status":"active","session_limit_minutes":45,"session_minutes_used":10}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	var resp contractsapi.CreatePlaybackSessionResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.PlaybackSessionID == "" || resp.Token == "" || resp.StreamURL == "" {
		t.Fatalf("expected playback session contract fields")
	}
}

func TestPostCreateSessionContractEntitlementRequired(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/playback/sessions",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"tablet","entitlement_status":"inactive"}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}

	var apiErr contractsapi.APIError
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("unmarshal api error: %v", err)
	}
	if apiErr.Code != "entitlement_required" {
		t.Fatalf("expected entitlement_required, got %q", apiErr.Code)
	}
}

func TestPostCreateSessionContractParentGateRequired(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/playback/sessions",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"external-link","entitlement_status":"active"}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}

	var apiErr contractsapi.APIError
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("unmarshal api error: %v", err)
	}
	if apiErr.Code != "parent_gate_required" {
		t.Fatalf("expected parent_gate_required, got %q", apiErr.Code)
	}
}

func TestPostCreateSessionContractSessionCapReached(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/playback/sessions",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"tablet","entitlement_status":"active","session_limit_minutes":15,"session_minutes_used":15}`),
	)
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}

	var apiErr contractsapi.APIError
	if err := json.Unmarshal(rr.Body.Bytes(), &apiErr); err != nil {
		t.Fatalf("unmarshal api error: %v", err)
	}
	if apiErr.Code != "session_cap_reached" {
		t.Fatalf("expected session_cap_reached, got %q", apiErr.Code)
	}
}

package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlaybackRouterEnforceMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	service := NewService()
	mux := NewMux(service)

	reqUnauthorized := httptest.NewRequest(http.MethodPost, "/v1/playback/sessions", strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"tablet"}`))
	reqUnauthorized.Header.Set("content-type", "application/json")
	rrUnauthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrUnauthorized, reqUnauthorized)
	if rrUnauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without identity headers, got %d", rrUnauthorized.Code)
	}

	reqAuthorized := httptest.NewRequest(http.MethodPost, "/v1/playback/sessions", strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","device_type":"tablet"}`))
	reqAuthorized.Header.Set("content-type", "application/json")
	reqAuthorized.Header.Set("X-Auth-Role", "child")
	rrAuthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrAuthorized, reqAuthorized)
	if rrAuthorized.Code != http.StatusCreated {
		t.Fatalf("expected 201 with child role, got %d", rrAuthorized.Code)
	}
}

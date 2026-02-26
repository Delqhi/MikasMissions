package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestGetKidsProgressContractHasCanonicalFields(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	watchReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/progress/watch-events",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","watch_ms":120000,"event_time":"2026-03-01T10:00:00Z"}`),
	)
	watchReq.Header.Set("content-type", "application/json")
	watchRR := httptest.NewRecorder()
	mux.ServeHTTP(watchRR, watchReq)
	if watchRR.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", watchRR.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/v1/kids/progress/child-1", nil)
	getRR := httptest.NewRecorder()
	mux.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRR.Code)
	}

	var resp contractsapi.KidsProgressResponse
	if err := json.Unmarshal(getRR.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.ChildProfileID == "" {
		t.Fatalf("expected child_profile_id")
	}
	if resp.WatchedMinutesToday < 0 || resp.WatchedMinutes7D < 0 {
		t.Fatalf("expected non-negative watch minute counters")
	}
	if resp.SessionLimitMinutes <= 0 {
		t.Fatalf("expected positive session limit")
	}
	if resp.LastEpisodeID != "ep-1" {
		t.Fatalf("expected last_episode_id ep-1, got %q", resp.LastEpisodeID)
	}
}

func TestPostWatchEventContractRejectsNegativeWatchMS(t *testing.T) {
	service := NewService()
	mux := NewMux(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/progress/watch-events",
		strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","watch_ms":-1,"event_time":"2026-03-01T10:00:00Z"}`),
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
	if apiErr.Code != "invalid_watch_ms" {
		t.Fatalf("expected invalid_watch_ms, got %q", apiErr.Code)
	}
}

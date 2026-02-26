package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetKidsProgress(t *testing.T) {
	service := NewService()
	mux := NewMux(service)
	watchReq := httptest.NewRequest(http.MethodPost, "/v1/progress/watch-events", strings.NewReader(`{"child_profile_id":"child-1","episode_id":"ep-1","watch_ms":2000,"event_time":"2026-03-01T10:00:00Z"}`))
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
}

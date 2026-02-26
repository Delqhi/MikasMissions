package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestParentControlsRoundTrip(t *testing.T) {
	store := NewStore()
	mux := NewMux(store)
	profileID := "child-1"
	body := strings.NewReader(`{"autoplay":false,"chat_enabled":false,"external_links":false,"session_limit_minutes":45,"bedtime_window":"20:00-07:00","safety_mode":"strict"}`)
	putReq := httptest.NewRequest(http.MethodPut, "/v1/parents/controls/"+profileID+"?parent_user_id=parent-1", body)
	putReq.Header.Set("content-type", "application/json")
	putRR := httptest.NewRecorder()
	mux.ServeHTTP(putRR, putReq)
	if putRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", putRR.Code)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/v1/parents/controls/"+profileID+"?parent_user_id=parent-1", nil)
	getRR := httptest.NewRecorder()
	mux.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRR.Code)
	}
	var resp contractsapi.ParentControlsResponse
	if err := json.Unmarshal(getRR.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Controls.SessionLimitMinutes != 45 {
		t.Fatalf("expected session limit 45, got %d", resp.Controls.SessionLimitMinutes)
	}
}

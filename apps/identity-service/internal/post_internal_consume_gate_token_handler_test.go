package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostInternalConsumeGateTokenOneTime(t *testing.T) {
	store := NewStore()
	if err := store.SaveGateToken("child-1", "token-1"); err != nil {
		t.Fatalf("save gate token: %v", err)
	}
	mux := NewMux(store)

	firstReq := httptest.NewRequest(http.MethodPost, "/internal/gates/consume", strings.NewReader(`{"child_profile_id":"child-1","gate_token":"token-1"}`))
	firstReq.Header.Set("content-type", "application/json")
	firstRR := httptest.NewRecorder()
	mux.ServeHTTP(firstRR, firstReq)
	if firstRR.Code != http.StatusOK || !strings.Contains(firstRR.Body.String(), `"consumed":true`) {
		t.Fatalf("expected first consume true, got %d %s", firstRR.Code, firstRR.Body.String())
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/internal/gates/consume", strings.NewReader(`{"child_profile_id":"child-1","gate_token":"token-1"}`))
	secondReq.Header.Set("content-type", "application/json")
	secondRR := httptest.NewRecorder()
	mux.ServeHTTP(secondRR, secondReq)
	if secondRR.Code != http.StatusOK || !strings.Contains(secondRR.Body.String(), `"consumed":false`) {
		t.Fatalf("expected second consume false, got %d %s", secondRR.Code, secondRR.Body.String())
	}
}

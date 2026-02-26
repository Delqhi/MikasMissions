package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProgressRouterEnforceMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	service := NewService()
	mux := NewMux(service)

	reqUnauthorized := httptest.NewRequest(http.MethodGet, "/v1/kids/progress/child-1", nil)
	rrUnauthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrUnauthorized, reqUnauthorized)
	if rrUnauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without identity headers, got %d", rrUnauthorized.Code)
	}

	reqAuthorized := httptest.NewRequest(http.MethodGet, "/v1/kids/progress/child-1", nil)
	reqAuthorized.Header.Set("X-Auth-Role", "child")
	rrAuthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrAuthorized, reqAuthorized)
	if rrAuthorized.Code != http.StatusOK {
		t.Fatalf("expected 200 with child role, got %d", rrAuthorized.Code)
	}
}

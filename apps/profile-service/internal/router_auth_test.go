package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProfileRouterEnforceMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	mux := NewMux(NewStore())

	reqUnauthorized := httptest.NewRequest(http.MethodPost, "/v1/children/profiles", strings.NewReader(`{"parent_user_id":"parent-1","display_name":"Mika","age_band":"6-11","avatar":"robot"}`))
	reqUnauthorized.Header.Set("content-type", "application/json")
	rrUnauthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrUnauthorized, reqUnauthorized)
	if rrUnauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without identity headers, got %d", rrUnauthorized.Code)
	}

	reqAuthorized := httptest.NewRequest(http.MethodPost, "/v1/children/profiles", strings.NewReader(`{"parent_user_id":"parent-1","display_name":"Mika","age_band":"6-11","avatar":"robot"}`))
	reqAuthorized.Header.Set("content-type", "application/json")
	reqAuthorized.Header.Set("X-Auth-Role", "parent")
	reqAuthorized.Header.Set("X-Auth-Parent-User-ID", "parent-1")
	rrAuthorized := httptest.NewRecorder()
	mux.ServeHTTP(rrAuthorized, reqAuthorized)
	if rrAuthorized.Code != http.StatusCreated {
		t.Fatalf("expected 201 with parent role, got %d", rrAuthorized.Code)
	}

	reqList := httptest.NewRequest(http.MethodGet, "/v1/children/profiles?parent_user_id=parent-1", nil)
	reqList.Header.Set("X-Auth-Role", "parent")
	reqList.Header.Set("X-Auth-Parent-User-ID", "parent-1")
	rrList := httptest.NewRecorder()
	mux.ServeHTTP(rrList, reqList)
	if rrList.Code != http.StatusOK {
		t.Fatalf("expected 200 for profile list with parent role, got %d", rrList.Code)
	}
}

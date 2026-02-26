package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEntitlementRequiresSubject(t *testing.T) {
	mux := NewMux(NewStore())
	req := httptest.NewRequest(http.MethodGet, "/v1/billing/entitlements", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetEntitlementByParent(t *testing.T) {
	mux := NewMux(NewStore())
	req := httptest.NewRequest(http.MethodGet, "/v1/billing/entitlements?parent_user_id=parent-1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetEntitlementByChild(t *testing.T) {
	mux := NewMux(NewStore())
	req := httptest.NewRequest(http.MethodGet, "/v1/billing/entitlements?child_profile_id=child-1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

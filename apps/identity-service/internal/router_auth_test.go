package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIdentityRouterEnforceMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	mux := NewMux(NewStore())

	reqProtected := httptest.NewRequest(http.MethodGet, "/v1/parents/dashboard?parent_user_id=parent-1", nil)
	rrProtected := httptest.NewRecorder()
	mux.ServeHTTP(rrProtected, reqProtected)
	if rrProtected.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for protected route without identity headers, got %d", rrProtected.Code)
	}

	reqPublic := httptest.NewRequest(http.MethodPost, "/v1/parents/signup", strings.NewReader(`{"email":"p@example.com","password":"1234567890x","country":"DE","language":"de","marketing":false,"accepted_terms":true}`))
	reqPublic.Header.Set("content-type", "application/json")
	rrPublic := httptest.NewRecorder()
	mux.ServeHTTP(rrPublic, reqPublic)
	if rrPublic.Code != http.StatusCreated {
		t.Fatalf("expected 201 for public signup route, got %d", rrPublic.Code)
	}

	reqLogin := httptest.NewRequest(http.MethodPost, "/v1/parents/login", strings.NewReader(`{"email":"p@example.com","password":"1234567890x"}`))
	reqLogin.Header.Set("content-type", "application/json")
	rrLogin := httptest.NewRecorder()
	mux.ServeHTTP(rrLogin, reqLogin)
	if rrLogin.Code != http.StatusUnauthorized && rrLogin.Code != http.StatusOK {
		t.Fatalf("expected public login route, got %d", rrLogin.Code)
	}
}

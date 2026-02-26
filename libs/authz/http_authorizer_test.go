package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPAuthorizerPermissiveWithoutIdentity(t *testing.T) {
	t.Setenv("AUTH_MODE", "permissive")
	authorizer := NewHTTPAuthorizerFromEnv()
	called := false
	handler := authorizer.Wrap([]string{"parent"}, func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !called {
		t.Fatalf("expected wrapped handler to be called")
	}
}

func TestHTTPAuthorizerEnforceWithoutIdentity(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	authorizer := NewHTTPAuthorizerFromEnv()
	handler := authorizer.Wrap([]string{"parent"}, func(_ http.ResponseWriter, _ *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestHTTPAuthorizerEnforceWithRole(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	authorizer := NewHTTPAuthorizerFromEnv()
	handler := authorizer.Wrap([]string{"parent"}, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set("X-Auth-Role", "parent")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestHTTPAuthorizerEnforceWrongRole(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	authorizer := NewHTTPAuthorizerFromEnv()
	handler := authorizer.Wrap([]string{"service"}, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set("X-Auth-Role", "parent")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPBillingEntitlementVerifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/billing/entitlements" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("child_profile_id"); got != "child-1" {
			t.Fatalf("unexpected child_profile_id: %s", got)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"active":true}`))
	}))
	defer server.Close()

	t.Setenv("BILLING_URL", server.URL)
	verifier, err := newHTTPBillingEntitlementVerifierFromEnv()
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	entitled, err := verifier.IsEntitled("child-1", "inactive")
	if err != nil {
		t.Fatalf("is entitled: %v", err)
	}
	if !entitled {
		t.Fatalf("expected entitled")
	}
}

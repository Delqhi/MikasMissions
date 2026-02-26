package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	return parsed
}

func newUpstream(name string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(name + ":" + r.Method + ":" + r.URL.Path))
	})
	return httptest.NewServer(handler)
}

func TestGatewayRoutesToExpectedUpstream(t *testing.T) {
	identity := newUpstream("identity")
	defer identity.Close()
	profile := newUpstream("profile")
	defer profile.Close()
	catalog := newUpstream("catalog")
	defer catalog.Close()
	recommendation := newUpstream("recommendation")
	defer recommendation.Close()
	playback := newUpstream("playback")
	defer playback.Close()
	progress := newUpstream("progress")
	defer progress.Close()
	creator := newUpstream("creator")
	defer creator.Close()
	billing := newUpstream("billing")
	defer billing.Close()
	adminStudio := newUpstream("admin-studio")
	defer adminStudio.Close()

	mux := NewMux(Upstreams{
		Identity:       mustParseURL(t, identity.URL),
		Profile:        mustParseURL(t, profile.URL),
		Catalog:        mustParseURL(t, catalog.URL),
		Recommendation: mustParseURL(t, recommendation.URL),
		Playback:       mustParseURL(t, playback.URL),
		Progress:       mustParseURL(t, progress.URL),
		Creator:        mustParseURL(t, creator.URL),
		Billing:        mustParseURL(t, billing.URL),
		AdminStudio:    mustParseURL(t, adminStudio.URL),
	})

	cases := []struct {
		method   string
		target   string
		body     string
		expected string
	}{
		{method: http.MethodPost, target: "/v1/parents/signup", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/parents/login", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/admin/login", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/parents/consent/verify", body: `{}`, expected: "identity"},
		{method: http.MethodGet, target: "/v1/parents/dashboard?parent_user_id=p1", expected: "identity"},
		{method: http.MethodGet, target: "/v1/parents/controls/child-1", expected: "identity"},
		{method: http.MethodPut, target: "/v1/parents/controls/child-1", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/parents/gates/challenge", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/parents/gates/verify", body: `{}`, expected: "identity"},
		{method: http.MethodPost, target: "/v1/children/profiles", body: `{}`, expected: "profile"},
		{method: http.MethodGet, target: "/v1/children/profiles?parent_user_id=p1", expected: "profile"},
		{method: http.MethodGet, target: "/v1/home/rails?child_profile_id=child-1", expected: "recommendation"},
		{method: http.MethodGet, target: "/v1/kids/home?child_profile_id=child-1&mode=core", expected: "recommendation"},
		{method: http.MethodGet, target: "/v1/kids/progress/child-1", expected: "progress"},
		{method: http.MethodGet, target: "/v1/catalog/episodes/ep-1", expected: "catalog"},
		{method: http.MethodPost, target: "/v1/playback/sessions", body: `{}`, expected: "playback"},
		{method: http.MethodPost, target: "/v1/progress/watch-events", body: `{}`, expected: "progress"},
		{method: http.MethodGet, target: "/v1/billing/entitlements?parent_user_id=p1", expected: "billing"},
		{method: http.MethodPost, target: "/v1/creator/assets/upload", body: `{}`, expected: "creator"},
		{method: http.MethodGet, target: "/v1/admin/workflows", expected: "admin-studio"},
		{method: http.MethodPost, target: "/v1/admin/workflows", body: `{}`, expected: "admin-studio"},
		{method: http.MethodPut, target: "/v1/admin/workflows/wf-1", body: `{}`, expected: "admin-studio"},
		{method: http.MethodDelete, target: "/v1/admin/workflows/wf-1", expected: "admin-studio"},
		{method: http.MethodPost, target: "/v1/admin/workflows/wf-1/runs", body: `{}`, expected: "admin-studio"},
		{method: http.MethodGet, target: "/v1/admin/runs/run-1", expected: "admin-studio"},
		{method: http.MethodGet, target: "/v1/admin/runs/run-1/logs", expected: "admin-studio"},
		{method: http.MethodPost, target: "/v1/admin/runs/run-1/retry", body: `{}`, expected: "admin-studio"},
		{method: http.MethodPost, target: "/v1/admin/runs/run-1/cancel", body: `{}`, expected: "admin-studio"},
		{method: http.MethodGet, target: "/v1/admin/model-profiles/default", expected: "admin-studio"},
		{method: http.MethodPut, target: "/v1/admin/model-profiles/default", body: `{}`, expected: "admin-studio"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s_%s", tc.method, tc.target), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, bytes.NewBufferString(tc.body))
			req.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rr.Code)
			}

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !bytes.HasPrefix(body, []byte(tc.expected+":")) {
				t.Fatalf("expected upstream %q, got body %q", tc.expected, string(body))
			}
		})
	}
}

func TestGatewayHealthz(t *testing.T) {
	mux := NewMux(Upstreams{})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("expected ok body, got %q", rr.Body.String())
	}
}

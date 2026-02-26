package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGatewayAuthEnforceMode(t *testing.T) {
	t.Setenv("AUTH_MODE", "enforce")
	t.Setenv("AUTH_JWT_SECRET", "test-secret")

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

	handler := NewMux(Upstreams{
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

	reqWithoutToken := httptest.NewRequest(http.MethodGet, "/v1/parents/dashboard?parent_user_id=p1", nil)
	rrWithoutToken := httptest.NewRecorder()
	handler.ServeHTTP(rrWithoutToken, reqWithoutToken)
	if rrWithoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rrWithoutToken.Code)
	}

	reqPublic := httptest.NewRequest(http.MethodPost, "/v1/parents/signup", nil)
	rrPublic := httptest.NewRecorder()
	handler.ServeHTTP(rrPublic, reqPublic)
	if rrPublic.Code != http.StatusOK {
		t.Fatalf("expected public signup to stay open, got %d", rrPublic.Code)
	}

	childToken := mustSignRoleToken(t, "test-secret", "child")
	reqWrongRole := httptest.NewRequest(http.MethodGet, "/v1/parents/dashboard?parent_user_id=p1", nil)
	reqWrongRole.Header.Set("Authorization", "Bearer "+childToken)
	rrWrongRole := httptest.NewRecorder()
	handler.ServeHTTP(rrWrongRole, reqWrongRole)
	if rrWrongRole.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong role, got %d", rrWrongRole.Code)
	}

	parentToken := mustSignRoleToken(t, "test-secret", "parent")
	reqParent := httptest.NewRequest(http.MethodGet, "/v1/parents/dashboard?parent_user_id=p1", nil)
	reqParent.Header.Set("Authorization", "Bearer "+parentToken)
	rrParent := httptest.NewRecorder()
	handler.ServeHTTP(rrParent, reqParent)
	if rrParent.Code != http.StatusOK {
		t.Fatalf("expected 200 for parent role, got %d", rrParent.Code)
	}

	reqCreator := httptest.NewRequest(http.MethodPost, "/v1/creator/assets/upload", nil)
	reqCreator.Header.Set("Authorization", "Bearer "+parentToken)
	rrCreator := httptest.NewRecorder()
	handler.ServeHTTP(rrCreator, reqCreator)
	if rrCreator.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for creator upload with parent role, got %d", rrCreator.Code)
	}

	reqAdminWithParent := httptest.NewRequest(http.MethodGet, "/v1/admin/workflows", nil)
	reqAdminWithParent.Header.Set("Authorization", "Bearer "+parentToken)
	rrAdminWithParent := httptest.NewRecorder()
	handler.ServeHTTP(rrAdminWithParent, reqAdminWithParent)
	if rrAdminWithParent.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for admin route with parent role, got %d", rrAdminWithParent.Code)
	}

	adminToken := mustSignRoleToken(t, "test-secret", "admin")
	reqAdmin := httptest.NewRequest(http.MethodGet, "/v1/admin/workflows", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminToken)
	rrAdmin := httptest.NewRecorder()
	handler.ServeHTTP(rrAdmin, reqAdmin)
	if rrAdmin.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin role, got %d", rrAdmin.Code)
	}
}

func TestGatewayAuthPermissiveModeAllowsLegacyClients(t *testing.T) {
	t.Setenv("AUTH_MODE", "permissive")
	t.Setenv("AUTH_JWT_SECRET", "")

	passThrough := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer passThrough.Close()
	parsed, err := url.Parse(passThrough.URL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	handler := NewMux(Upstreams{
		Identity:       parsed,
		Profile:        parsed,
		Catalog:        parsed,
		Recommendation: parsed,
		Playback:       parsed,
		Progress:       parsed,
		Creator:        parsed,
		Billing:        parsed,
		AdminStudio:    parsed,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/parents/dashboard?parent_user_id=p1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected permissive mode to allow without token, got %d", rr.Code)
	}
}

func mustSignRoleToken(t *testing.T, secret, role string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub":  "parent-test",
		"role": role,
		"exp":  time.Now().Add(10 * time.Minute).Unix(),
		"iat":  time.Now().Add(-1 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

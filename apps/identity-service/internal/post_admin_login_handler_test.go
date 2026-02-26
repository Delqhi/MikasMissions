package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/testkit"
)

func TestPostAdminLogin(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	store := NewStore()
	handler := PostAdminLogin(store)
	req := testkit.NewJSONRequest(http.MethodPost, "/v1/admin/login", contractsapi.AdminLoginRequest{
		Email:    "admin@mikasmissions.local",
		Password: "AdminPass123!",
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"role":"admin"`) {
		t.Fatalf("expected admin role in response")
	}
	if !strings.Contains(rr.Body.String(), `"access_token":"`) {
		t.Fatalf("expected access token in response")
	}
}

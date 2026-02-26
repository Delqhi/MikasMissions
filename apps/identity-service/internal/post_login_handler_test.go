package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/testkit"
)

func TestPostLogin(t *testing.T) {
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	store := NewStore()
	hashed, err := hashPassword("1234567890x")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	parent, err := store.CreateParent("parent@example.com", "DE", "de", hashed)
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	handler := PostLogin(store)
	req := testkit.NewJSONRequest(http.MethodPost, "/v1/parents/login", contractsapi.ParentLoginRequest{
		Email:    "parent@example.com",
		Password: "1234567890x",
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"parent_user_id":"`+parent.ID+`"`) {
		t.Fatalf("expected parent_user_id in response body")
	}
	if !strings.Contains(rr.Body.String(), `"access_token":"`) {
		t.Fatalf("expected access_token in response body")
	}
}

package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/testkit"
)

func TestPostSignup(t *testing.T) {
	store := NewStore()
	h := PostSignup(store)
	req := testkit.NewJSONRequest(http.MethodPost, "/v1/parents/signup", contractsapi.ParentSignupRequest{
		Email:      "parent@example.com",
		Password:   "1234567890x",
		Country:    "DE",
		Language:   "de",
		AcceptedTo: true,
	})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

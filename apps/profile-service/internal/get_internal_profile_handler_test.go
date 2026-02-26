package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetInternalProfile(t *testing.T) {
	store := NewStore()
	profile, err := store.CreateProfile("parent-1", "Mika", "6-11", "robot")
	if err != nil {
		t.Fatalf("create profile: %v", err)
	}
	mux := NewMux(store)

	req := httptest.NewRequest(http.MethodGet, "/internal/profiles/"+profile.ID, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, `"child_profile_id":"`+profile.ID+`"`) {
		t.Fatalf("expected child profile id in body, got %s", body)
	}
	if !strings.Contains(body, `"parent_user_id":"parent-1"`) {
		t.Fatalf("expected parent id in body, got %s", body)
	}
}

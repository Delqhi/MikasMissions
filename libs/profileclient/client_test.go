package profileclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProfileSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/profiles/child-1" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, `{"child_profile_id":"child-1","parent_user_id":"parent-1","age_band":"6-11"}`)
	}))
	defer server.Close()
	t.Setenv("PROFILE_URL", server.URL)
	client, err := NewFromEnv()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	profile, err := client.GetProfile(context.Background(), "child-1")
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if profile.ParentUserID != "parent-1" {
		t.Fatalf("expected parent-1, got %s", profile.ParentUserID)
	}
}

func TestGetProfileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	t.Setenv("PROFILE_URL", server.URL)
	client, err := NewFromEnv()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	_, err = client.GetProfile(context.Background(), "child-1")
	if err != ErrProfileNotFound {
		t.Fatalf("expected ErrProfileNotFound, got %v", err)
	}
}

func TestIsOwnedByParent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"child_profile_id":"child-1","parent_user_id":"parent-1","age_band":"6-11"}`)
	}))
	defer server.Close()
	t.Setenv("PROFILE_URL", server.URL)
	client, err := NewFromEnv()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	owned, err := client.IsOwnedByParent(context.Background(), "parent-1", "child-1")
	if err != nil {
		t.Fatalf("is owned: %v", err)
	}
	if !owned {
		t.Fatalf("expected owned")
	}
}

package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetKidsHome(t *testing.T) {
	service := NewService()
	handler := GetKidsHome(service)
	req := httptest.NewRequest(http.MethodGet, "/v1/kids/home?child_profile_id=child-1&mode=early", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

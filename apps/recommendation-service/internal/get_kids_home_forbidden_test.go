package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
)

type denyOwnerVerifier struct{}

func (denyOwnerVerifier) IsOwnedByParent(context.Context, string, string) (bool, error) {
	return false, nil
}

func TestGetKidsHomeForbiddenForForeignParent(t *testing.T) {
	service := NewServiceWithDependencies(newDemoCatalogReader(), denyOwnerVerifier{})
	handler := GetKidsHome(service)

	req := httptest.NewRequest(http.MethodGet, "/v1/kids/home?child_profile_id=child-1&mode=core", nil)
	req = req.WithContext(authz.WithPrincipal(req.Context(), authz.Principal{
		Role:         "parent",
		ParentUserID: "parent-foreign",
	}))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

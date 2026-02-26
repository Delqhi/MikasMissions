package internal

import (
	"context"
	"testing"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
)

type stubOwnerVerifier struct {
	owned bool
}

func (s stubOwnerVerifier) IsOwnedByParent(_ context.Context, _, _ string) (bool, error) {
	return s.owned, nil
}

func TestGetKidsProgressRejectsForeignParent(t *testing.T) {
	service := NewServiceWithRepositoryAndOwnerVerifier(NewStore(), stubOwnerVerifier{owned: false})
	ctx := authz.WithPrincipal(context.Background(), authz.Principal{
		ParentUserID: "parent-foreign",
		Role:         "parent",
	})
	_, err := service.GetKidsProgress(ctx, "child-1")
	if err != ErrChildProfileForbidden {
		t.Fatalf("expected ErrChildProfileForbidden, got %v", err)
	}
}

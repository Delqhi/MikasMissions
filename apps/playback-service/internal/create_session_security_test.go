package internal

import (
	"context"
	"testing"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type stubGateVerifier struct {
	consumed bool
	err      error
}

func (s stubGateVerifier) ConsumeToken(_ context.Context, _, _ string) (bool, error) {
	return s.consumed, s.err
}

type stubOwnerVerifier struct {
	owned bool
	err   error
}

func (s stubOwnerVerifier) IsOwnedByParent(_ context.Context, _, _ string) (bool, error) {
	return s.owned, s.err
}

func TestCreateSessionRejectsInvalidGateToken(t *testing.T) {
	service := NewServiceWithDependencies(
		stubEntitlementVerifier{entitled: true},
		stubGateVerifier{consumed: false},
		stubOwnerVerifier{owned: true},
	)
	_, err := service.CreateSession(context.Background(), contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID:  "child-1",
		EpisodeID:       "ep-1",
		DeviceType:      "external-link",
		ParentGateToken: "token-1",
		SafetyMode:      contractsapi.SafetyModeStrict,
	})
	if err != ErrParentGateRequired {
		t.Fatalf("expected ErrParentGateRequired, got %v", err)
	}
}

func TestCreateSessionRejectsForeignParentAccess(t *testing.T) {
	service := NewServiceWithDependencies(
		stubEntitlementVerifier{entitled: true},
		stubGateVerifier{consumed: true},
		stubOwnerVerifier{owned: false},
	)
	ctx := authz.WithPrincipal(context.Background(), authz.Principal{
		ParentUserID: "parent-foreign",
		Role:         "parent",
	})
	_, err := service.CreateSession(ctx, contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID: "child-1",
		EpisodeID:      "ep-1",
		DeviceType:     "tablet",
		SafetyMode:     contractsapi.SafetyModeStrict,
	})
	if err != ErrChildProfileForbidden {
		t.Fatalf("expected ErrChildProfileForbidden, got %v", err)
	}
}

package internal

import (
	"context"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type stubEntitlementVerifier struct {
	entitled bool
	err      error
}

func (s stubEntitlementVerifier) IsEntitled(_ string, _ string) (bool, error) {
	return s.entitled, s.err
}

func TestCreateSessionUsesVerifierResult(t *testing.T) {
	service := NewServiceWithEntitlementVerifier(stubEntitlementVerifier{entitled: false})
	_, err := service.CreateSession(context.Background(), contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID: "child-1",
		EpisodeID:      "ep-1",
		DeviceType:     "web",
		SafetyMode:     contractsapi.SafetyModeStrict,
	})
	if err == nil {
		t.Fatalf("expected entitlement error")
	}
}

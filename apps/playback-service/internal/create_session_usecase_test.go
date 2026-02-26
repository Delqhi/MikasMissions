package internal

import (
	"context"
	"testing"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func TestCreateSessionStrictAutoplayBlocked(t *testing.T) {
	service := NewService()
	_, err := service.CreateSession(context.Background(), contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID:    "child-1",
		EpisodeID:         "ep-1",
		DeviceType:        "web",
		SafetyMode:        contractsapi.SafetyModeStrict,
		EntitlementStatus: "active",
		AutoplayRequested: true,
	})
	if err == nil {
		t.Fatalf("expected error for autoplay in strict mode")
	}
}

func TestCreateSessionExternalLinkNeedsGate(t *testing.T) {
	service := NewService()
	_, err := service.CreateSession(context.Background(), contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID:    "child-1",
		EpisodeID:         "ep-1",
		DeviceType:        "external-link",
		SafetyMode:        contractsapi.SafetyModeStrict,
		EntitlementStatus: "active",
	})
	if err == nil {
		t.Fatalf("expected parent gate error")
	}
}

func TestCreateSessionEntitlementRequired(t *testing.T) {
	service := NewService()
	_, err := service.CreateSession(context.Background(), contractsapi.CreatePlaybackSessionRequest{
		ChildProfileID:    "child-1",
		EpisodeID:         "ep-1",
		DeviceType:        "web",
		SafetyMode:        contractsapi.SafetyModeStrict,
		EntitlementStatus: "inactive",
	})
	if err == nil {
		t.Fatalf("expected entitlement error")
	}
}

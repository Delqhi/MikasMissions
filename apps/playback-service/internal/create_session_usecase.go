package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
	"github.com/google/uuid"
)

type parentOwnershipVerifier interface {
	IsOwnedByParent(ctx context.Context, parentUserID, childProfileID string) (bool, error)
}

type Service struct {
	entitlementVerifier EntitlementVerifier
	gateVerifier        GateVerifier
	ownerVerifier       parentOwnershipVerifier
}

func NewService() *Service {
	return NewServiceWithDependencies(
		newDefaultEntitlementVerifier(),
		newDefaultGateVerifier(),
		profileclient.NewAllowAllReader(),
	)
}

func NewServiceWithEntitlementVerifier(entitlementVerifier EntitlementVerifier) *Service {
	return NewServiceWithDependencies(entitlementVerifier, newDefaultGateVerifier(), profileclient.NewAllowAllReader())
}

func NewServiceWithDependencies(
	entitlementVerifier EntitlementVerifier,
	gateVerifier GateVerifier,
	ownerVerifier parentOwnershipVerifier,
) *Service {
	if entitlementVerifier == nil {
		entitlementVerifier = newDefaultEntitlementVerifier()
	}
	if gateVerifier == nil {
		gateVerifier = newDefaultGateVerifier()
	}
	if ownerVerifier == nil {
		ownerVerifier = profileclient.NewAllowAllReader()
	}
	return &Service{
		entitlementVerifier: entitlementVerifier,
		gateVerifier:        gateVerifier,
		ownerVerifier:       ownerVerifier,
	}
}

func (s *Service) CreateSession(ctx context.Context, req contractsapi.CreatePlaybackSessionRequest) (contractsapi.CreatePlaybackSessionResponse, error) {
	if err := s.ensureParentOwnership(ctx, req.ChildProfileID); err != nil {
		return contractsapi.CreatePlaybackSessionResponse{}, err
	}
	if req.SafetyMode == "" {
		req.SafetyMode = contractsapi.SafetyModeStrict
	}
	sessionLimit := req.SessionLimitMinutes
	if sessionLimit == 0 {
		sessionLimit = 50
	}
	if sessionLimit > 120 {
		return contractsapi.CreatePlaybackSessionResponse{}, ErrSessionLimitTooHigh
	}
	if req.SessionMinutesUsed >= sessionLimit {
		return contractsapi.CreatePlaybackSessionResponse{}, ErrSessionCapReached
	}
	entitled, err := resolveEntitlementAllowed(s.entitlementVerifier, req)
	if err != nil {
		return contractsapi.CreatePlaybackSessionResponse{}, ErrEntitlementRequired
	}
	if !entitled {
		return contractsapi.CreatePlaybackSessionResponse{}, ErrEntitlementRequired
	}
	if req.SafetyMode == contractsapi.SafetyModeStrict && req.AutoplayRequested {
		return contractsapi.CreatePlaybackSessionResponse{}, ErrAutoplayBlocked
	}
	if req.DeviceType == "external-link" {
		if req.ParentGateToken == "" {
			return contractsapi.CreatePlaybackSessionResponse{}, ErrParentGateRequired
		}
		consumed, err := s.gateVerifier.ConsumeToken(ctx, req.ChildProfileID, req.ParentGateToken)
		if err != nil || !consumed {
			return contractsapi.CreatePlaybackSessionResponse{}, ErrParentGateRequired
		}
	}
	sessionID := uuid.NewString()
	token := uuid.NewString()
	sessionCapped := req.SessionMinutesUsed+15 >= sessionLimit
	return contractsapi.CreatePlaybackSessionResponse{
		PlaybackSessionID: sessionID,
		Token:             token,
		StreamURL:         fmt.Sprintf("https://stream.example.local/play/%s.m3u8", req.EpisodeID),
		ExpiresAt:         time.Now().UTC().Add(15 * time.Minute),
		SafetyApplied:     true,
		SafetyReason:      "strict_mode_active",
		SessionCapped:     sessionCapped,
		SessionMaxMinutes: sessionLimit,
	}, nil
}

func (s *Service) ensureParentOwnership(ctx context.Context, childProfileID string) error {
	principal, ok := authz.PrincipalFrom(ctx)
	if !ok || principal.Role != "parent" {
		return nil
	}
	if principal.ParentUserID == "" {
		return ErrChildProfileForbidden
	}
	owned, err := s.ownerVerifier.IsOwnedByParent(ctx, principal.ParentUserID, childProfileID)
	if err != nil {
		return err
	}
	if !owned {
		return ErrChildProfileForbidden
	}
	return nil
}

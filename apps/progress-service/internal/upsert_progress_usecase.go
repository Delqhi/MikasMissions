package internal

import (
	"context"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
)

type parentOwnershipVerifier interface {
	IsOwnedByParent(ctx context.Context, parentUserID, childProfileID string) (bool, error)
}

type Service struct {
	repository      Repository
	ownerVerifier   parentOwnershipVerifier
	defaultLimitMin int
}

func NewService() *Service {
	return NewServiceWithRepositoryAndOwnerVerifier(NewStore(), profileclient.NewAllowAllReader())
}

func NewServiceWithRepository(repository Repository) *Service {
	return NewServiceWithRepositoryAndOwnerVerifier(repository, profileclient.NewAllowAllReader())
}

func NewServiceWithRepositoryAndOwnerVerifier(repository Repository, ownerVerifier parentOwnershipVerifier) *Service {
	if ownerVerifier == nil {
		ownerVerifier = profileclient.NewAllowAllReader()
	}
	return &Service{
		repository:      repository,
		ownerVerifier:   ownerVerifier,
		defaultLimitMin: 50,
	}
}

func (s *Service) UpsertProgress(ctx context.Context, req contractsapi.UpsertWatchEventRequest) (contractsapi.UpsertWatchEventResponse, error) {
	if err := s.ensureParentOwnership(ctx, req.ChildProfileID); err != nil {
		return contractsapi.UpsertWatchEventResponse{}, err
	}
	eventTime := time.Now().UTC()
	if parsed, err := time.Parse(time.RFC3339, req.EventTimeISO); err == nil {
		eventTime = parsed.UTC()
	}
	if err := s.repository.AppendWatchEvent(ctx, req, eventTime); err != nil {
		return contractsapi.UpsertWatchEventResponse{}, err
	}
	return contractsapi.UpsertWatchEventResponse{Accepted: true}, nil
}

func (s *Service) GetKidsProgress(ctx context.Context, childProfileID string) (contractsapi.KidsProgressResponse, error) {
	if err := s.ensureParentOwnership(ctx, childProfileID); err != nil {
		return contractsapi.KidsProgressResponse{}, err
	}
	return s.repository.GetKidsProgress(ctx, childProfileID, time.Now().UTC(), s.defaultLimitMin)
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

package internal

import (
	"context"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
)

type catalogReader interface {
	ListEpisodes(ctx context.Context, ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error)
}

type parentOwnershipVerifier interface {
	IsOwnedByParent(ctx context.Context, parentUserID, childProfileID string) (bool, error)
}

type Service struct {
	catalogReader catalogReader
	ownerVerifier parentOwnershipVerifier
}

func NewService() *Service {
	return NewServiceWithDependencies(newCatalogReaderFromEnv(), profileclient.NewAllowAllReader())
}

func NewServiceWithDependencies(reader catalogReader, owner parentOwnershipVerifier) *Service {
	if reader == nil {
		reader = newCatalogReaderFromEnv()
	}
	if owner == nil {
		owner = profileclient.NewAllowAllReader()
	}
	return &Service{catalogReader: reader, ownerVerifier: owner}
}

func (s *Service) GetSafeRails(ctx context.Context, childProfileID string) (contractsapi.HomeRailsResponse, error) {
	if err := s.ensureParentOwnership(ctx, childProfileID); err != nil {
		return contractsapi.HomeRailsResponse{}, err
	}
	episodes, err := s.catalogReader.ListEpisodes(ctx, "", 20)
	if err != nil {
		return contractsapi.HomeRailsResponse{}, err
	}
	return contractsapi.HomeRailsResponse{
		ChildProfileID: childProfileID,
		Rails:          buildRails(episodes),
	}, nil
}

func (s *Service) GetKidsHome(ctx context.Context, childProfileID, mode string) (contractsapi.KidsHomeResponse, error) {
	if err := s.ensureParentOwnership(ctx, childProfileID); err != nil {
		return contractsapi.KidsHomeResponse{}, err
	}
	mode = normalizeMode(mode)
	episodes, err := s.catalogReader.ListEpisodes(ctx, ageBandForMode(mode), 20)
	if err != nil {
		return contractsapi.KidsHomeResponse{}, err
	}
	return contractsapi.KidsHomeResponse{
		ChildProfileID: childProfileID,
		Mode:           mode,
		SafetyMode:     contractsapi.SafetyModeStrict,
		PrimaryActions: actionsForMode(mode),
		Rails:          buildRails(episodes),
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

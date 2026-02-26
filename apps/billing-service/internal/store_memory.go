package internal

import (
	"context"
	"sync"
)

type Store struct {
	mu                sync.Mutex
	entitlementByUser map[string]EntitlementResponse
}

func NewStore() *Store {
	return &Store{
		entitlementByUser: map[string]EntitlementResponse{},
	}
}

func (s *Store) GetEntitlementByParent(_ context.Context, parentUserID string) (EntitlementResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.entitlementByUser[parentUserID]; ok {
		return existing, nil
	}
	defaultEntitlement := EntitlementResponse{
		ParentUserID: parentUserID,
		Plan:         "trial",
		Active:       true,
	}
	s.entitlementByUser[parentUserID] = defaultEntitlement
	return defaultEntitlement, nil
}

func (s *Store) GetEntitlementByChild(_ context.Context, childProfileID string) (EntitlementResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defaultEntitlement := EntitlementResponse{
		ParentUserID: "parent-" + childProfileID,
		Plan:         "trial",
		Active:       true,
	}
	return defaultEntitlement, nil
}

package internal

import (
	"sync"

	"github.com/google/uuid"
)

type Store struct {
	mu       sync.Mutex
	profiles map[string]ChildProfile
}

func NewStore() *Store {
	return &Store{profiles: map[string]ChildProfile{}}
}

func (s *Store) CreateProfile(parentID, displayName, ageBand, avatar string) (ChildProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile := ChildProfile{
		ID:          uuid.NewString(),
		ParentUser:  parentID,
		DisplayName: displayName,
		AgeBand:     ageBand,
		Avatar:      avatar,
	}
	s.profiles[profile.ID] = profile
	return profile, nil
}

func (s *Store) ListProfilesByParent(parentID string) ([]ChildProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profiles := make([]ChildProfile, 0, len(s.profiles))
	for _, profile := range s.profiles {
		if profile.ParentUser != parentID {
			continue
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (s *Store) FindProfile(id string) (ChildProfile, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile, ok := s.profiles[id]
	return profile, ok, nil
}

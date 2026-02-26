package internal

import (
	"sync"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/google/uuid"
)

type ParentGateChallenge struct {
	ChallengeID    string
	ParentUserID   string
	ChildProfileID string
	Method         string
	ExpiresAt      time.Time
	Used           bool
}

type gateTokenState struct {
	Token     string
	ExpiresAt time.Time
	Consumed  bool
}

type Store struct {
	mu            sync.Mutex
	parents       map[string]Parent
	admins        map[string]AdminUser
	consents      map[string]Consent
	controls      map[string]contractsapi.ParentalControls
	validGateByID map[string]gateTokenState
	challenges    map[string]ParentGateChallenge
}

func NewStore() *Store {
	return &Store{
		parents:       map[string]Parent{},
		admins:        seedDefaultAdmins(),
		consents:      map[string]Consent{},
		controls:      map[string]contractsapi.ParentalControls{},
		validGateByID: map[string]gateTokenState{},
		challenges:    map[string]ParentGateChallenge{},
	}
}

func (s *Store) CreateParent(email, country, lang, passwordHash string) (Parent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	parent := Parent{
		ID:           uuid.NewString(),
		Email:        email,
		Country:      country,
		Lang:         lang,
		PasswordHash: passwordHash,
	}
	s.parents[parent.ID] = parent
	return parent, nil
}

func (s *Store) FindParentByEmail(email string) (Parent, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, parent := range s.parents {
		if parent.Email == email {
			return parent, true, nil
		}
	}
	return Parent{}, false, nil
}

func (s *Store) UpdateParentLastLogin(_ string) error {
	return nil
}

func seedDefaultAdmins() map[string]AdminUser {
	result := map[string]AdminUser{}
	hash, err := hashPassword("AdminPass123!")
	if err != nil {
		return result
	}
	admin := AdminUser{
		ID:           "admin-local-1",
		Email:        "admin@mikasmissions.local",
		PasswordHash: hash,
	}
	result[admin.ID] = admin
	return result
}

func (s *Store) FindAdminByEmail(email string) (AdminUser, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, admin := range s.admins {
		if admin.Email == email {
			return admin, true, nil
		}
	}
	return AdminUser{}, false, nil
}

func (s *Store) UpdateAdminLastLogin(_ string) error {
	return nil
}

func (s *Store) VerifyConsent(parentID, method string) (Consent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	consent := Consent{
		ID:       uuid.NewString(),
		ParentID: parentID,
		Method:   method,
		Verified: true,
	}
	s.consents[consent.ID] = consent
	return consent, nil
}

func (s *Store) ParentExists(parentID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.parents[parentID]
	return ok, nil
}

func (s *Store) GetControls(childProfileID string) (contractsapi.ParentalControls, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	controls, ok := s.controls[childProfileID]
	if !ok {
		controls = contractsapi.DefaultStrictControls()
		s.controls[childProfileID] = controls
	}
	return controls, nil
}

func (s *Store) SetControls(childProfileID string, controls contractsapi.ParentalControls) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.controls[childProfileID] = controls
	return nil
}

func (s *Store) SaveGateToken(childProfileID, gateToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.validGateByID[childProfileID] = gateTokenState{
		Token:     gateToken,
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
		Consumed:  false,
	}
	return nil
}

func (s *Store) IsValidGateToken(childProfileID, gateToken string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.validGateByID[childProfileID]
	if !ok {
		return false, nil
	}
	if state.Consumed || time.Now().UTC().After(state.ExpiresAt) {
		return false, nil
	}
	return state.Token == gateToken && gateToken != "", nil
}

func (s *Store) ConsumeGateToken(childProfileID, gateToken string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.validGateByID[childProfileID]
	if !ok {
		return false, nil
	}
	if state.Consumed || time.Now().UTC().After(state.ExpiresAt) {
		return false, nil
	}
	if state.Token != gateToken || gateToken == "" {
		return false, nil
	}
	state.Consumed = true
	s.validGateByID[childProfileID] = state
	return true, nil
}

func (s *Store) CreateGateChallenge(parentUserID, childProfileID, method string) (ParentGateChallenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	challenge := ParentGateChallenge{
		ChallengeID:    uuid.NewString(),
		ParentUserID:   parentUserID,
		ChildProfileID: childProfileID,
		Method:         method,
		ExpiresAt:      time.Now().UTC().Add(5 * time.Minute),
	}
	s.challenges[challenge.ChallengeID] = challenge
	return challenge, nil
}

func (s *Store) ConsumeGateChallenge(challengeID, parentUserID, childProfileID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	challenge, ok := s.challenges[challengeID]
	if !ok {
		return false, nil
	}
	if challenge.Used || time.Now().UTC().After(challenge.ExpiresAt) {
		return false, nil
	}
	if challenge.ParentUserID != parentUserID || challenge.ChildProfileID != childProfileID {
		return false, nil
	}
	challenge.Used = true
	s.challenges[challengeID] = challenge
	return true, nil
}

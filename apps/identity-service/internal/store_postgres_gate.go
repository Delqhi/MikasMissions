package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *PostgresStore) CreateGateChallenge(parentUserID, childProfileID, method string) (ParentGateChallenge, error) {
	challenge := ParentGateChallenge{
		ChallengeID:    uuid.NewString(),
		ParentUserID:   parentUserID,
		ChildProfileID: childProfileID,
		Method:         method,
		ExpiresAt:      time.Now().UTC().Add(5 * time.Minute),
	}
	_, err := s.db.Exec(
		`insert into identity.parent_gate_challenges
		 (challenge_id, parent_user_id, child_profile_id, method, expires_at, used)
		 values ($1, $2, $3, $4, $5, false)`,
		challenge.ChallengeID,
		challenge.ParentUserID,
		challenge.ChildProfileID,
		challenge.Method,
		challenge.ExpiresAt,
	)
	if err != nil {
		return ParentGateChallenge{}, fmt.Errorf("insert gate challenge: %w", err)
	}
	return challenge, nil
}

func (s *PostgresStore) ConsumeGateChallenge(challengeID, parentUserID, childProfileID string) (bool, error) {
	var consumedID string
	err := s.db.QueryRow(
		`update identity.parent_gate_challenges
		 set used = true
		 where challenge_id = $1
		   and parent_user_id = $2
		   and child_profile_id = $3
		   and used = false
		   and expires_at > now()
		 returning challenge_id`,
		challengeID, parentUserID, childProfileID,
	).Scan(&consumedID)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("consume gate challenge: %w", err)
}

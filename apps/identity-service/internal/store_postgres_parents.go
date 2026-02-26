package internal

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (s *PostgresStore) CreateParent(email, country, lang, passwordHash string) (Parent, error) {
	parent := Parent{
		ID:           uuid.NewString(),
		Email:        email,
		Country:      country,
		Lang:         lang,
		PasswordHash: passwordHash,
	}
	_, err := s.db.Exec(
		`insert into identity.parents (id, email, country, language, password_hash) values ($1, $2, $3, $4, $5)`,
		parent.ID, parent.Email, parent.Country, parent.Lang, parent.PasswordHash,
	)
	if err != nil {
		return Parent{}, fmt.Errorf("insert parent: %w", err)
	}
	return parent, nil
}

func (s *PostgresStore) VerifyConsent(parentID, method string) (Consent, error) {
	consent := Consent{
		ID:       uuid.NewString(),
		ParentID: parentID,
		Method:   method,
		Verified: true,
	}
	_, err := s.db.Exec(
		`insert into identity.consents (id, parent_id, method, verified) values ($1, $2, $3, true)`,
		consent.ID, consent.ParentID, consent.Method,
	)
	if err != nil {
		return Consent{}, fmt.Errorf("insert consent: %w", err)
	}
	return consent, nil
}

func (s *PostgresStore) ParentExists(parentID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`select exists(select 1 from identity.parents where id::text = $1)`, parentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check parent exists: %w", err)
	}
	return exists, nil
}

func (s *PostgresStore) FindParentByEmail(email string) (Parent, bool, error) {
	var parent Parent
	err := s.db.QueryRow(
		`select id::text, email, country, language, coalesce(password_hash, '')
		 from identity.parents
		 where lower(email) = lower($1)`,
		email,
	).Scan(
		&parent.ID,
		&parent.Email,
		&parent.Country,
		&parent.Lang,
		&parent.PasswordHash,
	)
	if err == sql.ErrNoRows {
		return Parent{}, false, nil
	}
	if err != nil {
		return Parent{}, false, fmt.Errorf("find parent by email: %w", err)
	}
	return parent, true, nil
}

func (s *PostgresStore) UpdateParentLastLogin(parentID string) error {
	_, err := s.db.Exec(
		`update identity.parents
		 set last_login_at = now()
		 where id::text = $1`,
		parentID,
	)
	if err != nil {
		return fmt.Errorf("update parent last login: %w", err)
	}
	return nil
}

func (s *PostgresStore) IsValidGateToken(childProfileID, gateToken string) (bool, error) {
	var stored string
	err := s.db.QueryRow(
		`select gate_token
		 from identity.parent_gate_tokens
		 where child_profile_id = $1
		   and consumed_at is null
		   and expires_at > now()`,
		childProfileID,
	).Scan(&stored)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("fetch gate token: %w", err)
	}
	return stored == gateToken && gateToken != "", nil
}

package internal

import (
	"database/sql"
	"fmt"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func (s *PostgresStore) GetControls(childProfileID string) (contractsapi.ParentalControls, error) {
	var controls contractsapi.ParentalControls
	err := s.db.QueryRow(
		`select autoplay, chat_enabled, external_links, session_limit_minutes, bedtime_window, safety_mode
		 from identity.parent_controls
		 where child_profile_id = $1`,
		childProfileID,
	).Scan(
		&controls.Autoplay,
		&controls.ChatEnabled,
		&controls.ExternalLinks,
		&controls.SessionLimitMinutes,
		&controls.BedtimeWindow,
		&controls.SafetyMode,
	)
	if err == nil {
		return controls, nil
	}
	if err != sql.ErrNoRows {
		return contractsapi.ParentalControls{}, fmt.Errorf("query controls: %w", err)
	}
	defaultControls := contractsapi.DefaultStrictControls()
	if err := s.SetControls(childProfileID, defaultControls); err != nil {
		return contractsapi.ParentalControls{}, err
	}
	return defaultControls, nil
}

func (s *PostgresStore) SetControls(childProfileID string, controls contractsapi.ParentalControls) error {
	_, err := s.db.Exec(
		`insert into identity.parent_controls
		 (child_profile_id, autoplay, chat_enabled, external_links, session_limit_minutes, bedtime_window, safety_mode, updated_at)
		 values ($1, $2, $3, $4, $5, $6, $7, now())
		 on conflict (child_profile_id) do update set
		   autoplay = excluded.autoplay,
		   chat_enabled = excluded.chat_enabled,
		   external_links = excluded.external_links,
		   session_limit_minutes = excluded.session_limit_minutes,
		   bedtime_window = excluded.bedtime_window,
		   safety_mode = excluded.safety_mode,
		   updated_at = now()`,
		childProfileID,
		controls.Autoplay,
		controls.ChatEnabled,
		controls.ExternalLinks,
		controls.SessionLimitMinutes,
		controls.BedtimeWindow,
		controls.SafetyMode,
	)
	if err != nil {
		return fmt.Errorf("upsert controls: %w", err)
	}
	return nil
}

func (s *PostgresStore) SaveGateToken(childProfileID, gateToken string) error {
	_, err := s.db.Exec(
		`insert into identity.parent_gate_tokens (child_profile_id, gate_token, expires_at, consumed_at, updated_at)
		 values ($1, $2, now() + interval '5 minutes', null, now())
		 on conflict (child_profile_id) do update set
		   gate_token = excluded.gate_token,
		   expires_at = excluded.expires_at,
		   consumed_at = null,
		   updated_at = now()`,
		childProfileID, gateToken,
	)
	if err != nil {
		return fmt.Errorf("upsert gate token: %w", err)
	}
	return nil
}

func (s *PostgresStore) ConsumeGateToken(childProfileID, gateToken string) (bool, error) {
	var consumedID string
	err := s.db.QueryRow(
		`update identity.parent_gate_tokens
		 set consumed_at = now(),
		     updated_at = now()
		 where child_profile_id = $1
		   and gate_token = $2
		   and consumed_at is null
		   and expires_at > now()
		 returning child_profile_id`,
		childProfileID,
		gateToken,
	).Scan(&consumedID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("consume gate token: %w", err)
	}
	return true, nil
}

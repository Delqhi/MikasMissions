package internal

import (
	"database/sql"
	"fmt"
)

func (s *PostgresStore) GetModelProfile(modelProfileID string) (ModelProfile, bool, error) {
	var profile ModelProfile
	err := s.db.QueryRow(
		`select id, provider, base_url, model_id, timeout_ms, max_retries, safety_preset
		 from creator.model_profiles
		 where id = $1`,
		modelProfileID,
	).Scan(
		&profile.ID,
		&profile.Provider,
		&profile.BaseURL,
		&profile.ModelID,
		&profile.TimeoutMS,
		&profile.MaxRetries,
		&profile.SafetyPreset,
	)
	if err == sql.ErrNoRows {
		return ModelProfile{}, false, nil
	}
	if err != nil {
		return ModelProfile{}, false, fmt.Errorf("find model profile: %w", err)
	}
	return profile, true, nil
}

func (s *PostgresStore) PutModelProfile(profile ModelProfile, updatedBy string) (ModelProfile, error) {
	_, err := s.db.Exec(
		`insert into creator.model_profiles (id, provider, base_url, model_id, timeout_ms, max_retries, safety_preset, updated_at)
		 values ($1, $2, $3, $4, $5, $6, $7, now())
		 on conflict (id) do update set
		   provider = excluded.provider,
		   base_url = excluded.base_url,
		   model_id = excluded.model_id,
		   timeout_ms = excluded.timeout_ms,
		   max_retries = excluded.max_retries,
		   safety_preset = excluded.safety_preset,
		   updated_at = now()`,
		profile.ID,
		profile.Provider,
		profile.BaseURL,
		profile.ModelID,
		profile.TimeoutMS,
		profile.MaxRetries,
		profile.SafetyPreset,
	)
	if err != nil {
		return ModelProfile{}, fmt.Errorf("upsert model profile: %w", err)
	}
	if err := s.writeAuditAction(updatedBy, "model_profile_updated", "model_profile", profile.ID, nil); err != nil {
		return ModelProfile{}, err
	}
	return profile, nil
}

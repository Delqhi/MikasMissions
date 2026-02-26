package internal

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) CreateProfile(parentID, displayName, ageBand, avatar string) (ChildProfile, error) {
	profile := ChildProfile{
		ID:          uuid.NewString(),
		ParentUser:  parentID,
		DisplayName: displayName,
		AgeBand:     ageBand,
		Avatar:      avatar,
	}
	_, err := s.db.Exec(
		`insert into profiles.child_profiles (id, parent_id, display_name, age_band, avatar)
		 values ($1, $2, $3, $4, $5)`,
		profile.ID, profile.ParentUser, profile.DisplayName, profile.AgeBand, profile.Avatar,
	)
	if err != nil {
		return ChildProfile{}, fmt.Errorf("insert child profile: %w", err)
	}
	return profile, nil
}

func (s *PostgresStore) FindProfile(id string) (ChildProfile, bool, error) {
	var profile ChildProfile
	err := s.db.QueryRow(
		`select id::text, parent_id::text, display_name, age_band, avatar
		 from profiles.child_profiles
		 where id::text = $1`,
		id,
	).Scan(
		&profile.ID,
		&profile.ParentUser,
		&profile.DisplayName,
		&profile.AgeBand,
		&profile.Avatar,
	)
	if err == sql.ErrNoRows {
		return ChildProfile{}, false, nil
	}
	if err != nil {
		return ChildProfile{}, false, fmt.Errorf("find child profile: %w", err)
	}
	return profile, true, nil
}

func (s *PostgresStore) ListProfilesByParent(parentID string) ([]ChildProfile, error) {
	rows, err := s.db.Query(
		`select id::text, parent_id::text, display_name, age_band, avatar
		 from profiles.child_profiles
		 where parent_id::text = $1
		 order by created_at asc`,
		parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("list profiles by parent: %w", err)
	}
	defer rows.Close()
	profiles := make([]ChildProfile, 0)
	for rows.Next() {
		var profile ChildProfile
		if err := rows.Scan(
			&profile.ID,
			&profile.ParentUser,
			&profile.DisplayName,
			&profile.AgeBand,
			&profile.Avatar,
		); err != nil {
			return nil, fmt.Errorf("scan profile: %w", err)
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate profiles: %w", err)
	}
	return profiles, nil
}

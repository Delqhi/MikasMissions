package internal

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/delqhi/mikasmissions/platform/libs/generatorprovider"
	"github.com/delqhi/mikasmissions/platform/libs/runtimecfg"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type modelProfileReader interface {
	GetProfile(modelProfileID string) (generatorprovider.ModelProfile, error)
}

type staticModelProfileReader struct{}

type postgresModelProfileReader struct {
	db *sql.DB
}

func newModelProfileReaderFromEnv() modelProfileReader {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		store, err := newPostgresModelProfileReader(databaseURL)
		if err == nil {
			return store
		}
		if runtimecfg.PersistentStorageRequired() {
			panic(err)
		}
	}
	if runtimecfg.PersistentStorageRequired() {
		panic("DATABASE_URL is required for worker-gen-nim in strict persistence mode")
	}
	return staticModelProfileReader{}
}

func newPostgresModelProfileReader(databaseURL string) (*postgresModelProfileReader, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &postgresModelProfileReader{db: db}, nil
}

func (s staticModelProfileReader) GetProfile(_ string) (generatorprovider.ModelProfile, error) {
	return generatorprovider.ModelProfile{
		Provider:     "nvidia_nim",
		BaseURL:      envOr("NIM_BASE_URL", "http://127.0.0.1:9000"),
		ModelID:      envOr("NIM_MODEL_ID", "nim-video-v1"),
		TimeoutMS:    envOrInt("NIM_TIMEOUT_MS", 15000),
		MaxRetries:   envOrInt("NIM_MAX_RETRIES", 2),
		SafetyPreset: envOr("NIM_SAFETY_PRESET", "kids_strict"),
	}, nil
}

func (s *postgresModelProfileReader) GetProfile(modelProfileID string) (generatorprovider.ModelProfile, error) {
	if modelProfileID == "" {
		modelProfileID = "nim-default"
	}
	var profile generatorprovider.ModelProfile
	err := s.db.QueryRow(
		`select provider, base_url, model_id, timeout_ms, max_retries, safety_preset
		 from creator.model_profiles
		 where id = $1`,
		modelProfileID,
	).Scan(
		&profile.Provider,
		&profile.BaseURL,
		&profile.ModelID,
		&profile.TimeoutMS,
		&profile.MaxRetries,
		&profile.SafetyPreset,
	)
	if err != nil {
		return generatorprovider.ModelProfile{}, fmt.Errorf("query model profile: %w", err)
	}
	return profile, nil
}

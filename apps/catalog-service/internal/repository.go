package internal

import contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"

type Repository interface {
	FindEpisode(id string) (contractsapi.CatalogEpisodeResponse, bool, error)
	ListEpisodes(ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error)
	UpsertEpisode(req EpisodeUpsertRequest) error
}

type EpisodeUpsertRequest struct {
	EpisodeID      string   `json:"episode_id"`
	ShowID         string   `json:"show_id"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	AgeBand        string   `json:"age_band"`
	DurationMS     int64    `json:"duration_ms"`
	LearningTags   []string `json:"learning_tags"`
	PlaybackReady  bool     `json:"playback_ready"`
	ThumbnailURL   string   `json:"thumbnail_url"`
	PublishedAtISO string   `json:"published_at_iso"`
}

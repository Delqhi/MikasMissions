package internal

import (
	"database/sql"
	"fmt"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func (s *PostgresStore) FindEpisode(id string) (contractsapi.CatalogEpisodeResponse, bool, error) {
	var episode contractsapi.CatalogEpisodeResponse
	err := s.db.QueryRow(
		`select id, show_id, title, summary, age_band, duration_ms, learning_tags, playback_ready,
		        thumbnail_url, coalesce(published_at, now())::text
		 from catalog.episodes
		 where id = $1`,
		id,
	).Scan(
		&episode.EpisodeID,
		&episode.ShowID,
		&episode.Title,
		&episode.Summary,
		&episode.AgeBand,
		&episode.DurationMS,
		&episode.LearningTags,
		&episode.PlaybackReady,
		&episode.ThumbnailURL,
		&episode.PublishedAtISO,
	)
	if err == sql.ErrNoRows {
		return contractsapi.CatalogEpisodeResponse{}, false, nil
	}
	if err != nil {
		return contractsapi.CatalogEpisodeResponse{}, false, fmt.Errorf("find episode: %w", err)
	}
	return episode, true, nil
}

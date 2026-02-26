package internal

import (
	"fmt"
	"strings"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func (s *PostgresStore) ListEpisodes(ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `select id, show_id, title, summary, age_band, duration_ms, learning_tags, playback_ready,
	                 thumbnail_url, coalesce(published_at, now())::text
	          from catalog.episodes`
	args := []any{}
	if ageBand != "" {
		query += " where age_band = $1"
		args = append(args, ageBand)
	}
	query += " order by published_at desc nulls last, id desc"
	query += fmt.Sprintf(" limit %d", limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list episodes: %w", err)
	}
	defer rows.Close()

	episodes := make([]contractsapi.CatalogEpisodeResponse, 0, limit)
	for rows.Next() {
		var episode contractsapi.CatalogEpisodeResponse
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("scan listed episode: %w", err)
		}
		episode.PublishedAtISO = strings.ReplaceAll(episode.PublishedAtISO, " ", "T")
		episodes = append(episodes, episode)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listed episodes: %w", err)
	}
	return episodes, nil
}

package internal

import "fmt"

func (s *PostgresStore) UpsertEpisode(req EpisodeUpsertRequest) error {
	episode := applyEpisodeDefaults(req)
	_, err := s.db.Exec(
		`insert into catalog.episodes
		   (id, show_id, title, summary, age_band, duration_ms, learning_tags, playback_ready, thumbnail_url, published_at)
		 values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::timestamptz)
		 on conflict (id) do update set
		   show_id = excluded.show_id,
		   title = excluded.title,
		   summary = excluded.summary,
		   age_band = excluded.age_band,
		   duration_ms = excluded.duration_ms,
		   learning_tags = excluded.learning_tags,
		   playback_ready = excluded.playback_ready,
		   thumbnail_url = excluded.thumbnail_url,
		   published_at = excluded.published_at`,
		episode.EpisodeID,
		episode.ShowID,
		episode.Title,
		episode.Summary,
		episode.AgeBand,
		episode.DurationMS,
		episode.LearningTags,
		episode.PlaybackReady,
		episode.ThumbnailURL,
		episode.PublishedAtISO,
	)
	if err != nil {
		return fmt.Errorf("upsert episode: %w", err)
	}
	return nil
}

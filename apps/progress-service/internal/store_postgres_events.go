package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func (s *PostgresStore) AppendWatchEvent(ctx context.Context, req contractsapi.UpsertWatchEventRequest, eventTime time.Time) error {
	_, err := s.db.ExecContext(
		ctx,
		`insert into progress.watch_events (child_profile_id, episode_id, watch_ms, event_time)
		 values ($1, $2, $3, $4)`,
		req.ChildProfileID, req.EpisodeID, req.WatchMS, eventTime,
	)
	if err != nil {
		return fmt.Errorf("insert watch event: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetKidsProgress(
	ctx context.Context,
	childProfileID string,
	now time.Time,
	defaultLimitMin int,
) (contractsapi.KidsProgressResponse, error) {
	dayCutoff := now.Add(-24 * time.Hour)
	weekCutoff := now.Add(-7 * 24 * time.Hour)

	dayMS, weekMS, err := s.loadWindowSums(ctx, childProfileID, dayCutoff, weekCutoff)
	if err != nil {
		return contractsapi.KidsProgressResponse{}, err
	}
	lastEpisodeID, err := s.loadLastEpisodeID(ctx, childProfileID)
	if err != nil {
		return contractsapi.KidsProgressResponse{}, err
	}
	sessionLimitMin, err := s.loadSessionLimit(ctx, childProfileID, defaultLimitMin)
	if err != nil {
		return contractsapi.KidsProgressResponse{}, err
	}

	usedTodayMin := int(dayMS / 60000)
	completion := usedTodayMin * 2
	if completion > 100 {
		completion = 100
	}
	streak := 0
	if usedTodayMin > 0 {
		streak = 1
		if weekMS > dayMS {
			streak = 2
		}
	}

	return contractsapi.KidsProgressResponse{
		ChildProfileID:      childProfileID,
		WatchedMinutesToday: usedTodayMin,
		WatchedMinutes7D:    int(weekMS / 60000),
		CompletionPercent:   completion,
		MissionStreakDays:   streak,
		SessionLimitMinutes: sessionLimitMin,
		SessionMinutesUsed:  usedTodayMin,
		SessionCapped:       usedTodayMin >= sessionLimitMin,
		LastEpisodeID:       lastEpisodeID,
	}, nil
}

func (s *PostgresStore) loadWindowSums(
	ctx context.Context,
	childProfileID string,
	dayCutoff time.Time,
	weekCutoff time.Time,
) (int64, int64, error) {
	var dayMS int64
	var weekMS int64
	err := s.db.QueryRowContext(
		ctx,
		`select
		   coalesce(sum(case when event_time > $2 then watch_ms else 0 end), 0),
		   coalesce(sum(case when event_time > $3 then watch_ms else 0 end), 0)
		 from progress.watch_events
		 where child_profile_id::text = $1`,
		childProfileID, dayCutoff, weekCutoff,
	).Scan(&dayMS, &weekMS)
	if err != nil {
		return 0, 0, fmt.Errorf("query watch windows: %w", err)
	}
	return dayMS, weekMS, nil
}

func (s *PostgresStore) loadLastEpisodeID(ctx context.Context, childProfileID string) (string, error) {
	var lastEpisodeID string
	err := s.db.QueryRowContext(
		ctx,
		`select episode_id
		 from progress.watch_events
		 where child_profile_id::text = $1
		 order by event_time desc, id desc
		 limit 1`,
		childProfileID,
	).Scan(&lastEpisodeID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("query last episode: %w", err)
	}
	return lastEpisodeID, nil
}

func (s *PostgresStore) loadSessionLimit(ctx context.Context, childProfileID string, defaultLimitMin int) (int, error) {
	var sessionLimit int
	err := s.db.QueryRowContext(
		ctx,
		`select session_limit_minutes
		 from identity.parent_controls
		 where child_profile_id = $1`,
		childProfileID,
	).Scan(&sessionLimit)
	if err == sql.ErrNoRows {
		return defaultLimitMin, nil
	}
	if err != nil {
		return 0, fmt.Errorf("query session limit: %w", err)
	}
	return sessionLimit, nil
}

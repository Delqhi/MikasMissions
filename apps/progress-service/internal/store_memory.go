package internal

import (
	"context"
	"sync"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type watchRecord struct {
	episodeID string
	watchMS   int64
	eventTime time.Time
}

type Store struct {
	mu             sync.Mutex
	recordsByChild map[string][]watchRecord
}

func NewStore() *Store {
	return &Store{
		recordsByChild: map[string][]watchRecord{},
	}
}

func (s *Store) AppendWatchEvent(_ context.Context, req contractsapi.UpsertWatchEventRequest, eventTime time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recordsByChild[req.ChildProfileID] = append(s.recordsByChild[req.ChildProfileID], watchRecord{
		episodeID: req.EpisodeID,
		watchMS:   req.WatchMS,
		eventTime: eventTime,
	})
	return nil
}

func (s *Store) GetKidsProgress(_ context.Context, childProfileID string, now time.Time, defaultLimitMin int) (contractsapi.KidsProgressResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dayCutoff := now.Add(-24 * time.Hour)
	weekCutoff := now.Add(-7 * 24 * time.Hour)
	var dayMS int64
	var weekMS int64
	lastEpisode := ""
	var latestTime time.Time
	for _, rec := range s.recordsByChild[childProfileID] {
		if rec.eventTime.After(weekCutoff) {
			weekMS += rec.watchMS
		}
		if rec.eventTime.After(dayCutoff) {
			dayMS += rec.watchMS
		}
		if rec.eventTime.After(latestTime) {
			latestTime = rec.eventTime
			lastEpisode = rec.episodeID
		}
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
		SessionLimitMinutes: defaultLimitMin,
		SessionMinutesUsed:  usedTodayMin,
		SessionCapped:       usedTodayMin >= defaultLimitMin,
		LastEpisodeID:       lastEpisode,
	}, nil
}

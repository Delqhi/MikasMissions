package internal

import (
	"sort"
	"sync"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type Store struct {
	mu       sync.Mutex
	episodes map[string]contractsapi.CatalogEpisodeResponse
}

func NewStore() *Store {
	store := &Store{
		episodes: map[string]contractsapi.CatalogEpisodeResponse{},
	}
	_ = store.UpsertEpisode(EpisodeUpsertRequest{
		EpisodeID:      "ep-demo-1",
		ShowID:         "show-mikas-1",
		Title:          "Die Farben-Diebe",
		Summary:        "Mika und Team retten den Park.",
		AgeBand:        "6-11",
		DurationMS:     660000,
		LearningTags:   []string{"farben", "teamwork"},
		PlaybackReady:  true,
		ThumbnailURL:   "https://cdn.example.local/thumbs/ep-demo-1.jpg",
		PublishedAtISO: "2026-03-10T10:00:00Z",
	})
	return store
}

func (s *Store) FindEpisode(id string) (contractsapi.CatalogEpisodeResponse, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	episode, ok := s.episodes[id]
	return episode, ok, nil
}

func (s *Store) ListEpisodes(ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	episodes := make([]contractsapi.CatalogEpisodeResponse, 0, len(s.episodes))
	for _, episode := range s.episodes {
		if ageBand != "" && episode.AgeBand != ageBand {
			continue
		}
		episodes = append(episodes, episode)
	}
	sort.Slice(episodes, func(i, j int) bool {
		return publishedAt(episodes[i]).After(publishedAt(episodes[j]))
	})
	if limit > 0 && len(episodes) > limit {
		episodes = episodes[:limit]
	}
	return episodes, nil
}

func (s *Store) UpsertEpisode(req EpisodeUpsertRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	episode := applyEpisodeDefaults(req)
	s.episodes[episode.EpisodeID] = episode
	return nil
}

func publishedAt(episode contractsapi.CatalogEpisodeResponse) time.Time {
	parsed, err := time.Parse(time.RFC3339, episode.PublishedAtISO)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

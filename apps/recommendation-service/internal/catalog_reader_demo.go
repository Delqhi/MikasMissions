package internal

import (
	"context"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type demoCatalogReader struct{}

func newDemoCatalogReader() catalogReader {
	return &demoCatalogReader{}
}

func (r *demoCatalogReader) ListEpisodes(_ context.Context, ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error) {
	episodes := []contractsapi.CatalogEpisodeResponse{
		{
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
		},
		{
			EpisodeID:      "ep-demo-2",
			ShowID:         "show-mikas-1",
			Title:          "Sternen-Werkstatt",
			Summary:        "Ein Team-Abenteuer mit einfachen Missionen.",
			AgeBand:        "12-16",
			DurationMS:     780000,
			LearningTags:   []string{"wissenschaft", "kreativitaet"},
			PlaybackReady:  true,
			ThumbnailURL:   "https://cdn.example.local/thumbs/ep-demo-2.jpg",
			PublishedAtISO: "2026-03-11T10:00:00Z",
		},
		{
			EpisodeID:      "ep-demo-3",
			ShowID:         "show-mikas-1",
			Title:          "Regenbogen-Mission",
			Summary:        "Gefuehrte Entdeckungsreise fuer die Juengsten.",
			AgeBand:        "3-5",
			DurationMS:     420000,
			LearningTags:   []string{"farben", "audio-guided"},
			PlaybackReady:  true,
			ThumbnailURL:   "https://cdn.example.local/thumbs/ep-demo-3.jpg",
			PublishedAtISO: "2026-03-12T10:00:00Z",
		},
	}
	return filterEpisodes(episodes, ageBand, limit), nil
}

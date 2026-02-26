package internal

import (
	"strings"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func buildRails(episodes []contractsapi.CatalogEpisodeResponse) []contractsapi.RailItem {
	rails := make([]contractsapi.RailItem, 0, len(episodes))
	for _, episode := range episodes {
		suitability := "core"
		switch episode.AgeBand {
		case contractsapi.AgeBandEarly:
			suitability = "early"
		case contractsapi.AgeBandTeen:
			suitability = "teen"
		}
		rails = append(rails, contractsapi.RailItem{
			EpisodeID:          episode.EpisodeID,
			Title:              episode.Title,
			Summary:            episode.Summary,
			ThumbnailURL:       episode.ThumbnailURL,
			DurationMS:         episode.DurationMS,
			AgeBand:            episode.AgeBand,
			ContentSuitability: suitability,
			LearningTags:       episode.LearningTags,
			ReasonCode:         reasonCodeForEpisode(episode),
			SafetyApplied:      true,
			AgeFitScore:        0.93,
		})
	}
	return rails
}

func reasonCodeForEpisode(episode contractsapi.CatalogEpisodeResponse) string {
	if strings.Contains(strings.ToLower(episode.Title), "mission") {
		return "learning_path"
	}
	return "safe_curation"
}

func filterEpisodes(episodes []contractsapi.CatalogEpisodeResponse, ageBand string, limit int) []contractsapi.CatalogEpisodeResponse {
	filtered := make([]contractsapi.CatalogEpisodeResponse, 0, len(episodes))
	for _, episode := range episodes {
		if ageBand != "" && episode.AgeBand != ageBand {
			continue
		}
		filtered = append(filtered, episode)
		if limit > 0 && len(filtered) >= limit {
			break
		}
	}
	return filtered
}

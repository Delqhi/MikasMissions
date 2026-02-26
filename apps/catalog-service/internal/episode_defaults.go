package internal

import (
	"fmt"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func applyEpisodeDefaults(req EpisodeUpsertRequest) contractsapi.CatalogEpisodeResponse {
	showID := req.ShowID
	if showID == "" {
		showID = "show-mikas-1"
	}
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Mikas Mission %s", req.EpisodeID)
	}
	summary := req.Summary
	if summary == "" {
		summary = "Sicher kuratierte Episode aus der Publish-Pipeline."
	}
	durationMS := req.DurationMS
	if durationMS <= 0 {
		durationMS = 600000
	}
	thumbnailURL := req.ThumbnailURL
	if thumbnailURL == "" {
		thumbnailURL = fmt.Sprintf("https://cdn.example.local/thumbs/%s.jpg", req.EpisodeID)
	}
	publishedAtISO := req.PublishedAtISO
	if publishedAtISO == "" {
		publishedAtISO = time.Now().UTC().Format(time.RFC3339)
	}
	learningTags := req.LearningTags
	if len(learningTags) == 0 {
		learningTags = []string{"abenteuer"}
	}
	ageBand := req.AgeBand
	if !contractsapi.IsValidAgeBand(ageBand) {
		ageBand = contractsapi.AgeBandCore
	}
	return contractsapi.CatalogEpisodeResponse{
		EpisodeID:      req.EpisodeID,
		ShowID:         showID,
		Title:          title,
		Summary:        summary,
		AgeBand:        ageBand,
		DurationMS:     durationMS,
		LearningTags:   learningTags,
		PlaybackReady:  req.PlaybackReady,
		ThumbnailURL:   thumbnailURL,
		PublishedAtISO: publishedAtISO,
	}
}

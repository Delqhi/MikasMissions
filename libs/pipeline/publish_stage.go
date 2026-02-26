package pipeline

import (
	"strings"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
)

func BuildEpisodePublished(input contractsevents.MediaApprovedV1) (contractsevents.EpisodePublishedV1, error) {
	if err := input.Validate(); err != nil {
		return contractsevents.EpisodePublishedV1{}, err
	}
	episodeID := "ep-" + strings.ReplaceAll(input.AssetID, "-", "")
	if len(episodeID) > 14 {
		episodeID = episodeID[:14]
	}
	return contractsevents.EpisodePublishedV1{
		EpisodeID:    episodeID,
		AgeBand:      input.AgeBand,
		LearningTags: input.LearningTags,
	}, nil
}

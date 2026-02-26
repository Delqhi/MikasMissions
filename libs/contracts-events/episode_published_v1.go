package contractsevents

import "errors"

type EpisodePublishedV1 struct {
	EpisodeID    string   `json:"episode_id"`
	AgeBand      string   `json:"age_band"`
	LearningTags []string `json:"learning_tags"`
}

func (e EpisodePublishedV1) Validate() error {
	if e.EpisodeID == "" || e.AgeBand == "" {
		return errors.New("episode.published.v1 has missing required fields")
	}
	return nil
}

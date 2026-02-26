package contractsevents

import "errors"

type PlaybackSessionStartedV1 struct {
	PlaybackSessionID string `json:"playback_session_id"`
	ChildProfileID    string `json:"child_profile_id"`
	EpisodeID         string `json:"episode_id"`
	StartedAt         string `json:"started_at"`
	SafetyMode        string `json:"safety_mode"`
}

func (e PlaybackSessionStartedV1) Validate() error {
	if e.PlaybackSessionID == "" || e.ChildProfileID == "" || e.EpisodeID == "" || e.StartedAt == "" || e.SafetyMode == "" {
		return errors.New("playback.session.started.v1 has missing required fields")
	}
	return nil
}

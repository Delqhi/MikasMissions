package contractsevents

import "errors"

type PlaybackSessionEndedV1 struct {
	PlaybackSessionID string `json:"playback_session_id"`
	ChildProfileID    string `json:"child_profile_id"`
	EpisodeID         string `json:"episode_id"`
	EndedAt           string `json:"ended_at"`
	WatchedMS         int64  `json:"watched_ms"`
	Capped            bool   `json:"capped"`
}

func (e PlaybackSessionEndedV1) Validate() error {
	if e.PlaybackSessionID == "" || e.ChildProfileID == "" || e.EpisodeID == "" || e.EndedAt == "" {
		return errors.New("playback.session.ended.v1 has missing required fields")
	}
	if e.WatchedMS < 0 {
		return errors.New("playback.session.ended.v1 requires watched_ms >= 0")
	}
	return nil
}

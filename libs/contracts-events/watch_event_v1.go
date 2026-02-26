package contractsevents

import "errors"

type WatchEventV1 struct {
	ChildProfileID string `json:"child_profile_id"`
	EpisodeID      string `json:"episode_id"`
	WatchMS        int64  `json:"watch_ms"`
	EventTime      string `json:"event_time"`
}

func (e WatchEventV1) Validate() error {
	if e.ChildProfileID == "" || e.EpisodeID == "" || e.EventTime == "" {
		return errors.New("watch.event.v1 has missing required fields")
	}
	if e.WatchMS < 0 {
		return errors.New("watch.event.v1 watch_ms must be >= 0")
	}
	return nil
}

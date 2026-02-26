package contractsapi

type UpsertWatchEventRequest struct {
	ChildProfileID string `json:"child_profile_id"`
	EpisodeID      string `json:"episode_id"`
	WatchMS        int64  `json:"watch_ms"`
	EventTimeISO   string `json:"event_time"`
}

type UpsertWatchEventResponse struct {
	Accepted bool `json:"accepted"`
}

func (r UpsertWatchEventRequest) Validate() *APIError {
	if r.ChildProfileID == "" || r.EpisodeID == "" || r.EventTimeISO == "" {
		return &APIError{Code: "missing_fields", Message: "child_profile_id, episode_id and event_time are required"}
	}
	if r.WatchMS < 0 {
		return &APIError{Code: "invalid_watch_ms", Message: "watch_ms must be >= 0"}
	}
	return nil
}

package contractsapi

type KidsProgressResponse struct {
	ChildProfileID      string `json:"child_profile_id"`
	WatchedMinutesToday int    `json:"watched_minutes_today"`
	WatchedMinutes7D    int    `json:"watched_minutes_7d"`
	CompletionPercent   int    `json:"completion_percent"`
	MissionStreakDays   int    `json:"mission_streak_days"`
	SessionLimitMinutes int    `json:"session_limit_minutes"`
	SessionMinutesUsed  int    `json:"session_minutes_used"`
	SessionCapped       bool   `json:"session_capped"`
	LastEpisodeID       string `json:"last_episode_id"`
}

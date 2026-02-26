package contractsapi

type ParentsDashboardResponse struct {
	ParentUserID   string `json:"parent_user_id"`
	ChildProfiles  int    `json:"child_profiles"`
	WeeklyWatchMS  int64  `json:"weekly_watch_ms"`
	WeeklyLearning int    `json:"weekly_learning_events"`
}

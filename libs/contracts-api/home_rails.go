package contractsapi

type RailItem struct {
	EpisodeID          string   `json:"episode_id"`
	Title              string   `json:"title"`
	Summary            string   `json:"summary"`
	ThumbnailURL       string   `json:"thumbnail_url"`
	DurationMS         int64    `json:"duration_ms"`
	AgeBand            string   `json:"age_band"`
	ContentSuitability string   `json:"content_suitability"`
	LearningTags       []string `json:"learning_tags"`
	ReasonCode         string   `json:"reason_code"`
	SafetyApplied      bool     `json:"safety_applied"`
	AgeFitScore        float64  `json:"age_fit_score"`
}

type HomeRailsResponse struct {
	ChildProfileID string     `json:"child_profile_id"`
	Rails          []RailItem `json:"rails"`
}

type KidsHomeResponse struct {
	ChildProfileID string     `json:"child_profile_id"`
	Mode           string     `json:"mode"`
	SafetyMode     string     `json:"safety_mode"`
	PrimaryActions []string   `json:"primary_actions"`
	Rails          []RailItem `json:"rails"`
}

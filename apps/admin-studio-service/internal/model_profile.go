package internal

type ModelProfile struct {
	ID           string
	Provider     string
	BaseURL      string
	ModelID      string
	TimeoutMS    int
	MaxRetries   int
	SafetyPreset string
}

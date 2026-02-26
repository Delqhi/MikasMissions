package generatorprovider

import "encoding/json"

type ModelProfile struct {
	Provider     string
	BaseURL      string
	ModelID      string
	TimeoutMS    int
	MaxRetries   int
	SafetyPreset string
}

type GenerateRequest struct {
	RunID        string
	InputPayload json.RawMessage
}

type GenerateResult struct {
	AssetID    string
	SourceURL  string
	DurationMS int64
}

type Provider interface {
	GenerateVideo(req GenerateRequest) (GenerateResult, error)
}

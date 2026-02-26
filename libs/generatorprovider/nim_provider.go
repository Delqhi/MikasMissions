package generatorprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type nimProvider struct {
	client  *http.Client
	baseURL *url.URL
	modelID string
}

type nimRequest struct {
	ModelID      string          `json:"model_id"`
	InputPayload json.RawMessage `json:"input_payload"`
	RunID        string          `json:"run_id"`
}

type nimResponse struct {
	AssetID    string `json:"asset_id"`
	SourceURL  string `json:"source_url"`
	DurationMS int64  `json:"duration_ms"`
}

func NewNIMProvider(profile ModelProfile) Provider {
	timeout := profile.TimeoutMS
	if timeout <= 0 {
		timeout = 15000
	}
	parsed, _ := url.Parse(profile.BaseURL)
	if parsed == nil {
		parsed = &url.URL{Scheme: "http", Host: "127.0.0.1:9000"}
	}
	return &nimProvider{
		client:  &http.Client{Timeout: time.Duration(timeout) * time.Millisecond},
		baseURL: parsed,
		modelID: profile.ModelID,
	}
}

func (p *nimProvider) GenerateVideo(req GenerateRequest) (GenerateResult, error) {
	if p.baseURL == nil {
		return GenerateResult{}, fmt.Errorf("nim provider url is not configured")
	}
	requestURL := *p.baseURL
	requestURL.Path = "/v1/generate/video"
	payload, err := json.Marshal(nimRequest{
		ModelID:      p.modelID,
		InputPayload: req.InputPayload,
		RunID:        req.RunID,
	})
	if err != nil {
		return GenerateResult{}, fmt.Errorf("marshal nim request: %w", err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, requestURL.String(), bytes.NewReader(payload))
	if err != nil {
		return GenerateResult{}, fmt.Errorf("build nim request: %w", err)
	}
	httpReq.Header.Set("content-type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return GenerateResult{}, fmt.Errorf("nim request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return GenerateResult{}, fmt.Errorf("nim server error status: %d", resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return GenerateResult{}, fmt.Errorf("nim request rejected status: %d", resp.StatusCode)
	}
	var decoded nimResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return GenerateResult{}, fmt.Errorf("decode nim response: %w", err)
	}
	if decoded.AssetID == "" {
		decoded.AssetID = uuid.NewString()
	}
	if decoded.DurationMS <= 0 {
		decoded.DurationMS = 120000
	}
	if decoded.SourceURL == "" {
		decoded.SourceURL = fmt.Sprintf("https://cdn.generated.local/%s.mp4", decoded.AssetID)
	}
	return GenerateResult{
		AssetID:    decoded.AssetID,
		SourceURL:  decoded.SourceURL,
		DurationMS: decoded.DurationMS,
	}, nil
}

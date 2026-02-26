package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type httpCatalogProjector struct {
	client  *http.Client
	baseURL *url.URL
}

type noopCatalogProjector struct{}

type episodeProjectionRequest struct {
	EpisodeID      string   `json:"episode_id"`
	ShowID         string   `json:"show_id,omitempty"`
	Title          string   `json:"title,omitempty"`
	Summary        string   `json:"summary,omitempty"`
	AgeBand        string   `json:"age_band"`
	DurationMS     int64    `json:"duration_ms,omitempty"`
	LearningTags   []string `json:"learning_tags,omitempty"`
	PlaybackReady  bool     `json:"playback_ready"`
	ThumbnailURL   string   `json:"thumbnail_url,omitempty"`
	PublishedAtISO string   `json:"published_at_iso,omitempty"`
}

func NewCatalogProjectorFromEnv() catalogProjector {
	projector, err := newHTTPCatalogProjectorFromEnv()
	if err != nil {
		return noopCatalogProjector{}
	}
	return projector
}

func newHTTPCatalogProjectorFromEnv() (*httpCatalogProjector, error) {
	base := os.Getenv("CATALOG_URL")
	if base == "" {
		return nil, fmt.Errorf("CATALOG_URL is not set")
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse CATALOG_URL: %w", err)
	}
	return &httpCatalogProjector{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: parsed,
	}, nil
}

func (p *httpCatalogProjector) ProjectEpisode(ctx context.Context, req episodeProjectionRequest) error {
	requestURL := *p.baseURL
	requestURL.Path = "/internal/catalog/episodes"
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal upsert episode request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build upsert episode request: %w", err)
	}
	httpReq.Header.Set("content-type", "application/json")
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("project episode request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("project episode status code: %d", resp.StatusCode)
	}
	return nil
}

func (noopCatalogProjector) ProjectEpisode(context.Context, episodeProjectionRequest) error {
	return nil
}

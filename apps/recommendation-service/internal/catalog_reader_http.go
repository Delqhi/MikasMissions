package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type httpCatalogReader struct {
	client   *http.Client
	baseURL  *url.URL
	fallback catalogReader
	strict   bool
}

type listEpisodesResponse struct {
	Episodes []contractsapi.CatalogEpisodeResponse `json:"episodes"`
}

func newCatalogReaderFromEnv() catalogReader {
	reader, err := newHTTPCatalogReaderFromEnv()
	if err != nil {
		if strictRuntimeMode() {
			return failingCatalogReader{reason: err}
		}
		return newDemoCatalogReader()
	}
	return reader
}

func newHTTPCatalogReaderFromEnv() (*httpCatalogReader, error) {
	base := os.Getenv("CATALOG_URL")
	if base == "" {
		base = "http://127.0.0.1:8083"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse CATALOG_URL: %w", err)
	}
	return &httpCatalogReader{
		client:   &http.Client{Timeout: 2 * time.Second},
		baseURL:  parsed,
		fallback: newDemoCatalogReader(),
		strict:   strictRuntimeMode(),
	}, nil
}

func (r *httpCatalogReader) ListEpisodes(ctx context.Context, ageBand string, limit int) ([]contractsapi.CatalogEpisodeResponse, error) {
	requestURL := *r.baseURL
	requestURL.Path = "/internal/catalog/episodes"
	query := requestURL.Query()
	if ageBand != "" {
		query.Set("age_band", ageBand)
	}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	requestURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build list episodes request: %w", err)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		if r.strict {
			return nil, fmt.Errorf("request list episodes: %w", err)
		}
		return r.fallback.ListEpisodes(ctx, ageBand, limit)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if r.strict {
			return nil, fmt.Errorf("list episodes status code: %d", resp.StatusCode)
		}
		return r.fallback.ListEpisodes(ctx, ageBand, limit)
	}
	var decoded listEpisodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		if r.strict {
			return nil, fmt.Errorf("decode list episodes response: %w", err)
		}
		return r.fallback.ListEpisodes(ctx, ageBand, limit)
	}
	return filterEpisodes(decoded.Episodes, ageBand, limit), nil
}

type failingCatalogReader struct {
	reason error
}

func (f failingCatalogReader) ListEpisodes(context.Context, string, int) ([]contractsapi.CatalogEpisodeResponse, error) {
	return nil, f.reason
}

func strictRuntimeMode() bool {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("PERSISTENCE_MODE")))
	if mode == "strict" || mode == "required" {
		return true
	}
	env := strings.ToLower(strings.TrimSpace(os.Getenv("GO_ENV")))
	return env == "prod" || env == "production"
}

package profileclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Reader interface {
	GetProfile(ctx context.Context, childProfileID string) (Profile, error)
	IsOwnedByParent(ctx context.Context, parentUserID, childProfileID string) (bool, error)
}

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func NewFromEnv() (*Client, error) {
	base := os.Getenv("PROFILE_URL")
	if base == "" {
		base = "http://127.0.0.1:8082"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse PROFILE_URL: %w", err)
	}
	return &Client{
		httpClient: &http.Client{Timeout: 2 * time.Second},
		baseURL:    parsed,
	}, nil
}

func (c *Client) GetProfile(ctx context.Context, childProfileID string) (Profile, error) {
	if childProfileID == "" {
		return Profile{}, fmt.Errorf("child profile id is required")
	}
	requestURL := *c.baseURL
	requestURL.Path = "/internal/profiles/" + childProfileID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return Profile{}, fmt.Errorf("build profile request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("profile request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return Profile{}, ErrProfileNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("profile status code: %d", resp.StatusCode)
	}
	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return Profile{}, fmt.Errorf("decode profile response: %w", err)
	}
	return profile, nil
}

func (c *Client) IsOwnedByParent(ctx context.Context, parentUserID, childProfileID string) (bool, error) {
	profile, err := c.GetProfile(ctx, childProfileID)
	if err != nil {
		if err == ErrProfileNotFound {
			return false, nil
		}
		return false, err
	}
	return profile.ParentUserID == parentUserID, nil
}

type allowAllReader struct{}

func NewAllowAllReader() Reader {
	return allowAllReader{}
}

func (allowAllReader) GetProfile(_ context.Context, childProfileID string) (Profile, error) {
	return Profile{
		ChildProfileID: childProfileID,
		ParentUserID:   "parent-test",
		AgeBand:        "6-11",
	}, nil
}

func (allowAllReader) IsOwnedByParent(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

package internal

import (
	"fmt"
	"net/url"
	"os"
)

type Upstreams struct {
	Identity       *url.URL
	Profile        *url.URL
	Catalog        *url.URL
	Recommendation *url.URL
	Playback       *url.URL
	Progress       *url.URL
	Creator        *url.URL
	Billing        *url.URL
	AdminStudio    *url.URL
}

func loadURL(envKey, fallback string) (*url.URL, error) {
	raw := os.Getenv(envKey)
	if raw == "" {
		raw = fallback
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", envKey, err)
	}
	return parsed, nil
}

func LoadUpstreams() (Upstreams, error) {
	identity, err := loadURL("IDENTITY_URL", "http://127.0.0.1:8081")
	if err != nil {
		return Upstreams{}, err
	}
	profile, err := loadURL("PROFILE_URL", "http://127.0.0.1:8082")
	if err != nil {
		return Upstreams{}, err
	}
	catalog, err := loadURL("CATALOG_URL", "http://127.0.0.1:8083")
	if err != nil {
		return Upstreams{}, err
	}
	recommendation, err := loadURL("RECOMMENDATION_URL", "http://127.0.0.1:8084")
	if err != nil {
		return Upstreams{}, err
	}
	playback, err := loadURL("PLAYBACK_URL", "http://127.0.0.1:8085")
	if err != nil {
		return Upstreams{}, err
	}
	progress, err := loadURL("PROGRESS_URL", "http://127.0.0.1:8086")
	if err != nil {
		return Upstreams{}, err
	}
	creator, err := loadURL("CREATOR_URL", "http://127.0.0.1:8087")
	if err != nil {
		return Upstreams{}, err
	}
	billing, err := loadURL("BILLING_URL", "http://127.0.0.1:8089")
	if err != nil {
		return Upstreams{}, err
	}
	adminStudio, err := loadURL("ADMIN_STUDIO_URL", "http://127.0.0.1:8090")
	if err != nil {
		return Upstreams{}, err
	}
	return Upstreams{
		Identity:       identity,
		Profile:        profile,
		Catalog:        catalog,
		Recommendation: recommendation,
		Playback:       playback,
		Progress:       progress,
		Creator:        creator,
		Billing:        billing,
		AdminStudio:    adminStudio,
	}, nil
}

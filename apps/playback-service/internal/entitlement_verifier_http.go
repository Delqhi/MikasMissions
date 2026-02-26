package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type httpBillingEntitlementVerifier struct {
	client  *http.Client
	baseURL *url.URL
}

type billingEntitlementResponse struct {
	Active bool `json:"active"`
}

func newHTTPBillingEntitlementVerifierFromEnv() (*httpBillingEntitlementVerifier, error) {
	base := os.Getenv("BILLING_URL")
	if base == "" {
		base = "http://127.0.0.1:8089"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse BILLING_URL: %w", err)
	}
	return &httpBillingEntitlementVerifier{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: parsed,
	}, nil
}

func (v *httpBillingEntitlementVerifier) IsEntitled(childProfileID string, fallbackStatus string) (bool, error) {
	if childProfileID == "" {
		return false, fmt.Errorf("child_profile_id is required for entitlement lookup")
	}
	entitled, err := v.lookupChildEntitlement(childProfileID)
	if err != nil {
		return false, err
	}
	_ = fallbackStatus
	return entitled, nil
}

func (v *httpBillingEntitlementVerifier) lookupChildEntitlement(childProfileID string) (bool, error) {
	requestURL := *v.baseURL
	requestURL.Path = "/v1/billing/entitlements"
	query := requestURL.Query()
	query.Set("child_profile_id", childProfileID)
	requestURL.RawQuery = query.Encode()

	resp, err := v.client.Get(requestURL.String())
	if err != nil {
		return false, fmt.Errorf("billing entitlement request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("billing entitlement status code: %d", resp.StatusCode)
	}
	var decoded billingEntitlementResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return false, fmt.Errorf("decode billing entitlement: %w", err)
	}
	return decoded.Active, nil
}

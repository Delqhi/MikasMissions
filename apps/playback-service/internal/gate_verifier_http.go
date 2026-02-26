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

type httpGateVerifier struct {
	client  *http.Client
	baseURL *url.URL
}

type consumeGateTokenRequest struct {
	ChildProfileID string `json:"child_profile_id"`
	GateToken      string `json:"gate_token"`
}

type consumeGateTokenResponse struct {
	Consumed bool `json:"consumed"`
}

func NewIdentityGateVerifierFromEnv() GateVerifier {
	verifier, err := newHTTPGateVerifierFromEnv()
	if err != nil {
		return denyGateVerifier{reason: err}
	}
	return verifier
}

func newHTTPGateVerifierFromEnv() (*httpGateVerifier, error) {
	base := os.Getenv("IDENTITY_URL")
	if base == "" {
		base = "http://127.0.0.1:8081"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("parse IDENTITY_URL: %w", err)
	}
	return &httpGateVerifier{
		client:  &http.Client{Timeout: 2 * time.Second},
		baseURL: parsed,
	}, nil
}

func (v *httpGateVerifier) ConsumeToken(ctx context.Context, childProfileID, gateToken string) (bool, error) {
	requestURL := *v.baseURL
	requestURL.Path = "/internal/gates/consume"
	payload, err := json.Marshal(consumeGateTokenRequest{
		ChildProfileID: childProfileID,
		GateToken:      gateToken,
	})
	if err != nil {
		return false, fmt.Errorf("marshal consume gate token request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(payload))
	if err != nil {
		return false, fmt.Errorf("build consume gate token request: %w", err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := v.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("consume gate token request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("consume gate token status: %d", resp.StatusCode)
	}
	var decoded consumeGateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return false, fmt.Errorf("decode consume gate token response: %w", err)
	}
	return decoded.Consumed, nil
}

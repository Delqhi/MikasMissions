package internal

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type jwksProvider struct {
	jwksURL   string
	ttl       time.Duration
	client    *http.Client
	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	expiresAt time.Time
}

func newJWKSProvider(rawURL string, ttl time.Duration) (*jwksProvider, error) {
	if _, err := url.Parse(rawURL); err != nil {
		return nil, fmt.Errorf("parse jwks url: %w", err)
	}
	return &jwksProvider{
		jwksURL: rawURL,
		ttl:     ttl,
		client:  &http.Client{Timeout: 2 * time.Second},
		keys:    map[string]*rsa.PublicKey{},
	}, nil
}

func (p *jwksProvider) Key(kid string) (*rsa.PublicKey, error) {
	if kid == "" {
		return nil, fmt.Errorf("missing kid header")
	}
	p.mu.RLock()
	if key, ok := p.keys[kid]; ok && time.Now().Before(p.expiresAt) {
		p.mu.RUnlock()
		return key, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	if key, ok := p.keys[kid]; ok && time.Now().Before(p.expiresAt) {
		return key, nil
	}
	if err := p.refresh(); err != nil {
		return nil, err
	}
	key, ok := p.keys[kid]
	if !ok {
		return nil, fmt.Errorf("kid %s not found in jwks", kid)
	}
	return key, nil
}

func (p *jwksProvider) refresh() error {
	resp, err := p.client.Get(p.jwksURL)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch jwks status: %d", resp.StatusCode)
	}

	var payload struct {
		Keys []struct {
			KID string `json:"kid"`
			KTY string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}
	parsed := map[string]*rsa.PublicKey{}
	for _, key := range payload.Keys {
		if key.KTY != "RSA" || key.KID == "" {
			continue
		}
		publicKey, err := parseRSAPublicKey(key.N, key.E)
		if err != nil {
			continue
		}
		parsed[key.KID] = publicKey
	}
	if len(parsed) == 0 {
		return fmt.Errorf("no usable rsa keys in jwks")
	}
	p.keys = parsed
	p.expiresAt = time.Now().Add(p.ttl)
	return nil
}

func parseRSAPublicKey(rawModulus, rawExponent string) (*rsa.PublicKey, error) {
	modulusBytes, err := base64.RawURLEncoding.DecodeString(rawModulus)
	if err != nil {
		return nil, fmt.Errorf("decode modulus: %w", err)
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(rawExponent)
	if err != nil {
		return nil, fmt.Errorf("decode exponent: %w", err)
	}
	modulus := new(big.Int).SetBytes(modulusBytes)
	exponent := int(new(big.Int).SetBytes(exponentBytes).Int64())
	if modulus.Sign() <= 0 || exponent <= 0 {
		return nil, fmt.Errorf("invalid rsa key parameters")
	}
	return &rsa.PublicKey{N: modulus, E: exponent}, nil
}

package internal

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/golang-jwt/jwt/v5"
)

type tokenVerifier struct {
	issuer   string
	audience string
	secret   []byte
	jwks     *jwksProvider
}

func newTokenVerifierFromEnv() (*tokenVerifier, error) {
	secret := strings.TrimSpace(os.Getenv("AUTH_JWT_SECRET"))
	jwksURL := firstNonEmpty(
		strings.TrimSpace(os.Getenv("AUTH_JWKS_URL")),
		strings.TrimSpace(os.Getenv("SUPABASE_JWKS_URL")),
		defaultJWKSURLFromSupabaseURL(strings.TrimSpace(os.Getenv("SUPABASE_URL"))),
	)
	var jwksProvider *jwksProvider
	var err error
	if jwksURL != "" {
		jwksProvider, err = newJWKSProvider(jwksURL, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}
	if secret == "" && jwksProvider == nil {
		return nil, nil
	}
	return &tokenVerifier{
		issuer:   strings.TrimSpace(os.Getenv("AUTH_JWT_ISSUER")),
		audience: strings.TrimSpace(os.Getenv("AUTH_JWT_AUDIENCE")),
		secret:   []byte(secret),
		jwks:     jwksProvider,
	}, nil
}

func defaultJWKSURLFromSupabaseURL(rawSupabaseURL string) string {
	if rawSupabaseURL == "" {
		return ""
	}
	parsed, err := url.Parse(rawSupabaseURL)
	if err != nil {
		return ""
	}
	parsed.Path = "/auth/v1/.well-known/jwks.json"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func (v *tokenVerifier) Verify(rawToken string) (authz.Principal, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		v.keyFunc,
		jwt.WithValidMethods([]string{"HS256", "RS256"}),
	)
	if err != nil {
		return authz.Principal{}, fmt.Errorf("parse token: %w", err)
	}
	if !parsedToken.Valid {
		return authz.Principal{}, fmt.Errorf("invalid token")
	}
	if err := v.validateIssuer(claims); err != nil {
		return authz.Principal{}, err
	}
	if err := v.validateAudience(claims); err != nil {
		return authz.Principal{}, err
	}
	return principalFromClaims(claims)
}

func (v *tokenVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	switch token.Method.Alg() {
	case "HS256":
		if len(v.secret) == 0 {
			return nil, fmt.Errorf("hs256 token not supported without AUTH_JWT_SECRET")
		}
		return v.secret, nil
	case "RS256":
		if v.jwks == nil {
			return nil, fmt.Errorf("rs256 token not supported without AUTH_JWKS_URL")
		}
		kid, _ := token.Header["kid"].(string)
		return v.jwks.Key(kid)
	default:
		return nil, fmt.Errorf("unsupported jwt algorithm")
	}
}

func (v *tokenVerifier) validateIssuer(claims jwt.MapClaims) error {
	if v.issuer == "" {
		return nil
	}
	if claimsIssuer := stringClaim(claims, "iss"); claimsIssuer != v.issuer {
		return fmt.Errorf("invalid issuer")
	}
	return nil
}

func (v *tokenVerifier) validateAudience(claims jwt.MapClaims) error {
	if v.audience == "" {
		return nil
	}
	raw := claims["aud"]
	switch audience := raw.(type) {
	case string:
		if audience == v.audience {
			return nil
		}
	case []interface{}:
		for _, item := range audience {
			if value, ok := item.(string); ok && value == v.audience {
				return nil
			}
		}
	}
	return fmt.Errorf("invalid audience")
}

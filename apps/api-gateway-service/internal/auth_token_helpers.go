package internal

import (
	"errors"
	"strings"
)

var errMissingBearerToken = errors.New("missing bearer token")

func bearerToken(rawAuthorization string) (string, error) {
	if rawAuthorization == "" {
		return "", errMissingBearerToken
	}
	parts := strings.SplitN(rawAuthorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errMissingBearerToken
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errMissingBearerToken
	}
	return token, nil
}

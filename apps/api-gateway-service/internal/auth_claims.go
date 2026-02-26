package internal

import (
	"fmt"
	"strings"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/golang-jwt/jwt/v5"
)

func principalFromClaims(claims jwt.MapClaims) (authz.Principal, error) {
	role := firstNonEmpty(
		stringClaim(claims, "mm_role"),
		stringClaim(claims, "role"),
		nestedStringClaim(claims, "app_metadata", "mm_role"),
		nestedStringClaim(claims, "app_metadata", "role"),
		nestedStringClaim(claims, "user_metadata", "mm_role"),
		nestedStringClaim(claims, "user_metadata", "role"),
		stringClaim(claims, "user_role"),
	)
	if role == "" {
		return authz.Principal{}, fmt.Errorf("missing role claim")
	}
	parentUserID := firstNonEmpty(
		stringClaim(claims, "parent_user_id"),
		stringClaim(claims, "sub"),
	)
	return authz.Principal{
		ParentUserID: parentUserID,
		Role:         strings.ToLower(role),
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func stringClaim(claims jwt.MapClaims, key string) string {
	raw, ok := claims[key]
	if !ok {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func nestedStringClaim(claims jwt.MapClaims, parent, key string) string {
	rawParent, ok := claims[parent]
	if !ok {
		return ""
	}
	parentMap, ok := rawParent.(map[string]interface{})
	if !ok {
		return ""
	}
	rawValue, ok := parentMap[key]
	if !ok {
		return ""
	}
	value, ok := rawValue.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

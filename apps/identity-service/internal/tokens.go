package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	parentTokenTTL = 3600
	adminTokenTTL  = 3600
)

func issueParentToken(parentUserID string) (string, int, error) {
	return issueRoleToken(parentUserID, "parent", parentTokenTTL)
}

func issueAdminToken(adminUserID string) (string, int, error) {
	return issueRoleToken(adminUserID, "admin", adminTokenTTL)
}

func issueRoleToken(subjectID, role string, ttlSeconds int) (string, int, error) {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		return "", 0, fmt.Errorf("AUTH_JWT_SECRET is required")
	}
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  subjectID,
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Duration(ttlSeconds) * time.Second).Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, fmt.Errorf("sign role token: %w", err)
	}
	return signed, ttlSeconds, nil
}

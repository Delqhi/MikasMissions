package internal

import (
	"os"
	"strings"
)

const (
	authModePermissive = "permissive"
	authModeEnforce    = "enforce"
)

type authConfig struct {
	mode string
}

func loadAuthConfigFromEnv() authConfig {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_MODE")))
	if mode != authModeEnforce {
		mode = authModePermissive
	}
	return authConfig{mode: mode}
}

func (c authConfig) requireTokenForProtectedRoutes() bool {
	return c.mode == authModeEnforce
}

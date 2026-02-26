package runtimecfg

import (
	"os"
	"strings"
)

func PersistentStorageRequired() bool {
	return PersistentStorageRequiredFromEnv(os.Getenv)
}

func PersistentStorageRequiredFromEnv(getenv func(string) string) bool {
	mode := strings.ToLower(strings.TrimSpace(getenv("PERSISTENCE_MODE")))
	switch mode {
	case "strict", "required":
		return true
	}
	env := strings.ToLower(strings.TrimSpace(getenv("GO_ENV")))
	switch env {
	case "prod", "production":
		return true
	}
	return false
}

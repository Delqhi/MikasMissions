package authz

import (
	"net/http"
	"os"
	"strings"

	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

const (
	AuthModePermissive = "permissive"
	AuthModeEnforce    = "enforce"
)

type HTTPAuthorizer struct {
	mode string
}

func NewHTTPAuthorizerFromEnv() *HTTPAuthorizer {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_MODE")))
	if mode != AuthModeEnforce {
		mode = AuthModePermissive
	}
	return &HTTPAuthorizer{mode: mode}
}

func (a *HTTPAuthorizer) Wrap(allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(allowedRoles) == 0 {
			next(w, r)
			return
		}
		role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Auth-Role")))
		parentUserID := strings.TrimSpace(r.Header.Get("X-Auth-Parent-User-ID"))
		if role == "" {
			if a.mode == AuthModeEnforce {
				httpx.WriteAPIError(w, http.StatusUnauthorized, "missing_identity", "authenticated identity headers are required")
				return
			}
			next(w, r)
			return
		}
		if !containsRole(role, allowedRoles) {
			httpx.WriteAPIError(w, http.StatusForbidden, "insufficient_role", "role is not allowed for this operation")
			return
		}
		authorized := r.WithContext(WithPrincipal(r.Context(), Principal{
			ParentUserID: parentUserID,
			Role:         role,
		}))
		next(w, authorized)
	}
}

func containsRole(role string, allowed []string) bool {
	for _, candidate := range allowed {
		if role == candidate {
			return true
		}
	}
	return false
}

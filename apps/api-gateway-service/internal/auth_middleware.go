package internal

import (
	"log"
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

type authError struct {
	status  int
	code    string
	message string
}

type gatewayAuth struct {
	config   authConfig
	verifier *tokenVerifier
}

func newGatewayAuthFromEnv() *gatewayAuth {
	config := loadAuthConfigFromEnv()
	verifier, err := newTokenVerifierFromEnv()
	if err != nil {
		log.Printf("gateway auth verifier init failed: %v", err)
	}
	return &gatewayAuth{
		config:   config,
		verifier: verifier,
	}
}

func (a *gatewayAuth) Authorize(r *http.Request, pattern string) (*http.Request, *authError) {
	requiredRoles := requiredRolesForPattern(pattern)
	if len(requiredRoles) == 0 {
		return r, nil
	}

	rawToken, err := bearerToken(r.Header.Get("Authorization"))
	if err != nil {
		if a.config.requireTokenForProtectedRoutes() {
			return nil, &authError{
				status:  http.StatusUnauthorized,
				code:    "missing_token",
				message: "bearer token required",
			}
		}
		return r, nil
	}
	if a.verifier == nil {
		return nil, &authError{
			status:  http.StatusServiceUnavailable,
			code:    "auth_unavailable",
			message: "token verification unavailable",
		}
	}
	principal, verifyErr := a.verifier.Verify(rawToken)
	if verifyErr != nil {
		return nil, &authError{
			status:  http.StatusUnauthorized,
			code:    "invalid_token",
			message: verifyErr.Error(),
		}
	}
	if !roleAllowed(principal.Role, requiredRoles) {
		return nil, &authError{
			status:  http.StatusForbidden,
			code:    "insufficient_role",
			message: "token role is not allowed for this route",
		}
	}

	authorized := r.WithContext(authz.WithPrincipal(r.Context(), principal))
	authorized.Header.Set("X-Auth-Role", principal.Role)
	if principal.ParentUserID != "" {
		authorized.Header.Set("X-Auth-Parent-User-ID", principal.ParentUserID)
	}
	return authorized, nil
}

type gatewayHandler struct {
	mux  *http.ServeMux
	auth *gatewayAuth
}

func newGatewayHandler(mux *http.ServeMux, auth *gatewayAuth) http.Handler {
	return &gatewayHandler{mux: mux, auth: auth}
}

func (h *gatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, pattern := h.mux.Handler(r)
	if h.auth != nil && pattern != "" {
		authorizedRequest, authErr := h.auth.Authorize(r, pattern)
		if authErr != nil {
			httpx.WriteAPIError(w, authErr.status, authErr.code, authErr.message)
			return
		}
		r = authorizedRequest
	}
	handler.ServeHTTP(w, r)
}

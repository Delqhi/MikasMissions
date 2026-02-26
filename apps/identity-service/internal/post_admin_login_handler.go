package internal

import (
	"errors"
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"golang.org/x/crypto/bcrypt"
)

func PostAdminLogin(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.AdminLoginRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		admin, found, err := store.FindAdminByEmail(req.Email)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "login_error", err.Error())
			return
		}
		if !found || admin.PasswordHash == "" {
			httpx.WriteAPIError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is invalid")
			return
		}
		if err := verifyPassword(admin.PasswordHash, req.Password); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				httpx.WriteAPIError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is invalid")
				return
			}
			httpx.WriteAPIError(w, http.StatusInternalServerError, "login_error", err.Error())
			return
		}
		token, expiresIn, err := issueAdminToken(admin.ID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "login_error", err.Error())
			return
		}
		if err := store.UpdateAdminLastLogin(admin.ID); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "login_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.AdminLoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
			AdminUserID: admin.ID,
			Role:        "admin",
		})
	}
}

package internal

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
)

func PostSignup(store Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.ParentSignupRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "signup_error", err.Error())
			return
		}
		parent, err := store.CreateParent(req.Email, req.Country, req.Language, passwordHash)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "signup_error", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, contractsapi.ParentSignupResponse{
			ParentUserID: parent.ID,
			Status:       "pending_consent",
		})
	}
}

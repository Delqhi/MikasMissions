package internal

import (
	"net/http"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/httpx"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

func PostConsentVerify(store Repository, bus queue.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req contractsapi.ParentConsentVerifyRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}
		if apiErr := req.Validate(); apiErr != nil {
			httpx.WriteJSON(w, http.StatusBadRequest, apiErr)
			return
		}
		exists, err := store.ParentExists(req.ParentUserID)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "consent_error", err.Error())
			return
		}
		if !exists {
			httpx.WriteAPIError(w, http.StatusNotFound, "parent_not_found", "parent_user_id not found")
			return
		}
		consent, err := store.VerifyConsent(req.ParentUserID, req.Method)
		if err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "consent_error", err.Error())
			return
		}
		if err := publishEvent(r.Context(), bus, "consent.verified.v1", uuid.NewString(), contractsevents.ConsentVerifiedV1{
			ConsentID:    consent.ID,
			ParentUserID: req.ParentUserID,
			Method:       req.Method,
			VerifiedAt:   time.Now().UTC().Format(time.RFC3339),
		}); err != nil {
			httpx.WriteAPIError(w, http.StatusInternalServerError, "event_publish_failed", err.Error())
			return
		}
		httpx.WriteJSON(w, http.StatusOK, contractsapi.ParentConsentVerifyResponse{
			ConsentID: consent.ID,
			Verified:  consent.Verified,
		})
	}
}

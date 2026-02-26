package httpx

import (
	"net/http"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

func WriteAPIError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, contractsapi.APIError{Code: code, Message: message})
}

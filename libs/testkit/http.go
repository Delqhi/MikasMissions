package testkit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func NewJSONRequest(method, target string, payload any) *http.Request {
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

package contractsapi

import "time"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AuditMeta struct {
	TraceID   string    `json:"trace_id"`
	CreatedAt time.Time `json:"created_at"`
}

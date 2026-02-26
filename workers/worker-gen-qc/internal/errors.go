package internal

import "errors"

var (
	errInvalidSourceURL = errors.New("source_url must be absolute http(s) url")
	errInvalidDuration  = errors.New("duration_ms outside qc bounds")
)

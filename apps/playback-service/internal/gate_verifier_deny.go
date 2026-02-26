package internal

import (
	"context"
	"fmt"
)

type denyGateVerifier struct {
	reason error
}

func (v denyGateVerifier) ConsumeToken(_ context.Context, _ string, _ string) (bool, error) {
	if v.reason != nil {
		return false, fmt.Errorf("gate verifier unavailable: %w", v.reason)
	}
	return false, fmt.Errorf("gate verifier unavailable")
}

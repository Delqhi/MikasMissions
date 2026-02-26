package internal

import "fmt"

type denyEntitlementVerifier struct {
	reason error
}

func (v denyEntitlementVerifier) IsEntitled(_ string, _ string) (bool, error) {
	if v.reason != nil {
		return false, fmt.Errorf("entitlement verifier unavailable: %w", v.reason)
	}
	return false, fmt.Errorf("entitlement verifier unavailable")
}

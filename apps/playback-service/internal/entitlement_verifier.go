package internal

import contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"

type EntitlementVerifier interface {
	IsEntitled(childProfileID string, fallbackStatus string) (bool, error)
}

func resolveFallbackEntitlement(fallbackStatus string) bool {
	if fallbackStatus == "" {
		fallbackStatus = "active"
	}
	return fallbackStatus == "active"
}

func newDefaultEntitlementVerifier() EntitlementVerifier {
	return &fallbackEntitlementVerifier{}
}

func NewBillingEntitlementVerifierFromEnv() EntitlementVerifier {
	verifier, err := newHTTPBillingEntitlementVerifierFromEnv()
	if err != nil {
		return denyEntitlementVerifier{reason: err}
	}
	return verifier
}

func resolveEntitlementAllowed(verifier EntitlementVerifier, req contractsapi.CreatePlaybackSessionRequest) (bool, error) {
	if verifier == nil {
		verifier = newDefaultEntitlementVerifier()
	}
	return verifier.IsEntitled(req.ChildProfileID, req.EntitlementStatus)
}

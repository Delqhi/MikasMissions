package internal

type fallbackEntitlementVerifier struct{}

func (v *fallbackEntitlementVerifier) IsEntitled(_ string, fallbackStatus string) (bool, error) {
	return resolveFallbackEntitlement(fallbackStatus), nil
}

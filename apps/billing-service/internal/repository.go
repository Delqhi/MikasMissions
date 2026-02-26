package internal

import "context"

type Repository interface {
	GetEntitlementByParent(ctx context.Context, parentUserID string) (EntitlementResponse, error)
	GetEntitlementByChild(ctx context.Context, childProfileID string) (EntitlementResponse, error)
}

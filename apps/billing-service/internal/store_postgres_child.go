package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *PostgresStore) GetEntitlementByChild(ctx context.Context, childProfileID string) (EntitlementResponse, error) {
	var parentID string
	var plan string
	var active bool
	err := s.db.QueryRowContext(
		ctx,
		`select cp.parent_id::text,
		        coalesce(sub.plan, 'trial') as plan,
		        coalesce(sub.active, true) as active
		 from profiles.child_profiles cp
		 left join lateral (
		   select plan, active
		   from billing.subscriptions
		   where parent_id = cp.parent_id
		   order by created_at desc
		   limit 1
		 ) sub on true
		 where cp.id::text = $1`,
		childProfileID,
	).Scan(&parentID, &plan, &active)
	if errors.Is(err, sql.ErrNoRows) {
		return EntitlementResponse{}, ErrEntitlementNotFound
	}
	if err != nil {
		return EntitlementResponse{}, fmt.Errorf("query entitlement by child: %w", err)
	}
	return EntitlementResponse{
		ParentUserID: parentID,
		Plan:         plan,
		Active:       active,
	}, nil
}

package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *PostgresStore) GetEntitlementByParent(ctx context.Context, parentUserID string) (EntitlementResponse, error) {
	var parentID string
	var plan string
	var active bool
	err := s.db.QueryRowContext(
		ctx,
		`select p.id::text,
		        coalesce(sub.plan, 'trial') as plan,
		        coalesce(sub.active, true) as active
		 from identity.parents p
		 left join lateral (
		   select plan, active
		   from billing.subscriptions
		   where parent_id = p.id
		   order by created_at desc
		   limit 1
		 ) sub on true
		 where p.id::text = $1`,
		parentUserID,
	).Scan(&parentID, &plan, &active)
	if errors.Is(err, sql.ErrNoRows) {
		return EntitlementResponse{}, ErrEntitlementNotFound
	}
	if err != nil {
		return EntitlementResponse{}, fmt.Errorf("query entitlement by parent: %w", err)
	}
	return EntitlementResponse{
		ParentUserID: parentID,
		Plan:         plan,
		Active:       active,
	}, nil
}

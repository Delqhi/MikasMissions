package internal

import (
	"context"
	"time"

	contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"
)

type Repository interface {
	AppendWatchEvent(ctx context.Context, req contractsapi.UpsertWatchEventRequest, eventTime time.Time) error
	GetKidsProgress(ctx context.Context, childProfileID string, now time.Time, defaultLimitMin int) (contractsapi.KidsProgressResponse, error)
}

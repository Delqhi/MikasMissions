package contractsapi

import "context"

type CatalogService interface {
	ListRails(ctx context.Context, childProfileID string) (HomeRailsResponse, error)
}

type PlaybackService interface {
	CreateSession(ctx context.Context, req CreatePlaybackSessionRequest) (CreatePlaybackSessionResponse, error)
}

type ModerationService interface {
	EvaluateAsset(ctx context.Context, assetID string, ageBand string) (string, error)
}

type RecommendationService interface {
	GetSafeRails(ctx context.Context, childProfileID string) (HomeRailsResponse, error)
	GetKidsHome(ctx context.Context, childProfileID, mode string) (KidsHomeResponse, error)
}

type ProgressService interface {
	UpsertProgress(ctx context.Context, req UpsertWatchEventRequest) (UpsertWatchEventResponse, error)
	GetKidsProgress(ctx context.Context, childProfileID string) (KidsProgressResponse, error)
}

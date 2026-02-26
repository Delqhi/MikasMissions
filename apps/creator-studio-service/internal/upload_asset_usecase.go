package internal

import (
	"context"
	"encoding/json"
	"io"

	contractsevents "github.com/delqhi/mikasmissions/platform/libs/contracts-events"
	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

type UploadRequest struct {
	SourceURL  string `json:"source_url"`
	UploaderID string `json:"uploader_id"`
}

type UploadResponse struct {
	AssetID string `json:"asset_id"`
	Status  string `json:"status"`
}

type Service struct {
	bus          queue.Bus
	outboxWriter outboxWriter
	outboxCloser io.Closer
}

func NewService(bus queue.Bus) (*Service, error) {
	writer, closer, err := newOutboxWriterFromEnv()
	if err != nil {
		return nil, err
	}
	return &Service{bus: bus, outboxWriter: writer, outboxCloser: closer}, nil
}

func (s *Service) Close() error {
	if s.outboxCloser == nil {
		return nil
	}
	return s.outboxCloser.Close()
}

func (s *Service) UploadAsset(ctx context.Context, req UploadRequest) (UploadResponse, error) {
	assetID := uuid.NewString()
	eventBody, err := json.Marshal(contractsevents.MediaUploadedV1{
		AssetID:   assetID,
		SourceURL: req.SourceURL,
		Uploader:  req.UploaderID,
		TraceID:   uuid.NewString(),
	})
	if err != nil {
		return UploadResponse{}, err
	}
	if err := s.outboxWriter.EnqueueAndFlush(ctx, s.bus, queue.Event{
		ID:      uuid.NewString(),
		Topic:   "media.uploaded.v1",
		Payload: eventBody,
	}); err != nil {
		return UploadResponse{}, err
	}
	return UploadResponse{AssetID: assetID, Status: "queued"}, nil
}

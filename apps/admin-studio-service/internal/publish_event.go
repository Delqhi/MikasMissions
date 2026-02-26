package internal

import (
	"context"
	"encoding/json"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
	"github.com/google/uuid"
)

func publishJSONEvent(ctx context.Context, bus queue.Bus, topic string, payload any) error {
	if bus == nil {
		return nil
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return bus.Publish(ctx, queue.Event{
		ID:      uuid.NewString(),
		Topic:   topic,
		Payload: encoded,
	})
}

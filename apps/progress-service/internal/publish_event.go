package internal

import (
	"context"
	"encoding/json"

	"github.com/delqhi/mikasmissions/platform/libs/queue"
)

func publishEvent(ctx context.Context, bus queue.Bus, topic, eventID string, payload any) error {
	if bus == nil {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return bus.Publish(ctx, queue.Event{
		ID:      eventID,
		Topic:   topic,
		Payload: body,
	})
}

package queue

import (
	"os"
)

func NewBusFromEnv() (Bus, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return NewInMemoryBus(), nil
	}
	return NewNATSBus(natsURL)
}

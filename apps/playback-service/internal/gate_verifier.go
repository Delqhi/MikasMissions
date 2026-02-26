package internal

import "context"

type GateVerifier interface {
	ConsumeToken(ctx context.Context, childProfileID, gateToken string) (bool, error)
}

type allowAllGateVerifier struct{}

func (allowAllGateVerifier) ConsumeToken(_ context.Context, _ string, gateToken string) (bool, error) {
	return gateToken != "", nil
}

func newDefaultGateVerifier() GateVerifier {
	return allowAllGateVerifier{}
}

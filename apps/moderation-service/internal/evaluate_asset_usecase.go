package internal

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) EvaluateAsset(_ context.Context, _ string, ageBand string) (string, error) {
	if ageBand == "3-5" || ageBand == "6-11" || ageBand == "12-16" {
		return "approved", nil
	}
	return "needs_review", nil
}

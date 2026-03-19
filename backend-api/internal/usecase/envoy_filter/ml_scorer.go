package envoy_filter

import (
	"context"
	"errors"
)

var errMLScorerUnavailable = errors.New("ml scorer is not configured")

type UnavailableMLScorer struct{}

func NewUnavailableMLScorer() *UnavailableMLScorer {
	return &UnavailableMLScorer{}
}

func (s *UnavailableMLScorer) Score(_ context.Context, _ string, _ []byte, _ []float32) (float32, error) {
	return 0, errMLScorerUnavailable
}

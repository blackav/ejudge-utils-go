package gentypes

import (
	"context"
	"time"
)

type ModelInfo struct {
	ID      string
	Object  string
	OwnedBy string
	Type    string
}

type CompletionResultUsage struct {
	InputTokens  *int64
	OutputTokens *int64
	TotalTokens  *int64
}

type CompletionResult struct {
	Model     string
	CreatedAt time.Time
	Text      string
	Usage     *CompletionResultUsage
}

type Generator interface {
	RefreshToken(ctx context.Context) error
	ListModels(ctx context.Context) ([]ModelInfo, error)
	SimpleCompletion(ctx context.Context, model string, text string) (*CompletionResult, error)
}

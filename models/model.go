package models

import (
	"context"

	"github.com/oliver-platt/goagent/v2/types"
)

// Model represents an LLM that can generate responses
type Model interface {
	Generate(ctx context.Context, messages []types.Message) (string, error)
	Name() string
}

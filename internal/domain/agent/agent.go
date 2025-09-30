package agent

import (
	"context"
)

type AgentClient interface {
	Generate(ctx context.Context, prompt string) ([]byte, error)
}

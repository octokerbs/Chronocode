package agent

import (
	"context"
)

type GenerativeAgentService interface {
	Generate(ctx context.Context, prompt string) ([]byte, error)
}

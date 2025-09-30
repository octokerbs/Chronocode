package domain

import (
	"context"
)

type Agent interface {
	Generate(ctx context.Context, prompt string) ([]byte, error)
}

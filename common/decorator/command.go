package decorator

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *zap.Logger) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base:   handler,
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

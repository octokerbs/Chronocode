package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type GenerativeAgentServiceMock struct {
	mock.Mock
}

func (g *GenerativeAgentServiceMock) Generate(ctx context.Context, prompt string) ([]byte, error) {
	return nil, nil
}

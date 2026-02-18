package query

import (
	"context"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
)

type GetUserProfile struct {
	AccessToken string
}

type GetUserProfileHandler struct {
	codeHostFactory codehost.CodeHostFactory
}

func NewGetUserProfileHandler(codeHostFactory codehost.CodeHostFactory) GetUserProfileHandler {
	return GetUserProfileHandler{codeHostFactory: codeHostFactory}
}

func (h *GetUserProfileHandler) Handle(ctx context.Context, cmd GetUserProfile) (*codehost.UserProfile, error) {
	ch, err := h.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		return nil, err
	}

	return ch.GetAuthenticatedUser(ctx)
}

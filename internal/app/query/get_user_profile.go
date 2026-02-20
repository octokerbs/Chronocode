package query

import (
	"context"
	"log/slog"

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
	slog.Info("GetUserProfile query received")

	ch, err := h.codeHostFactory.Create(ctx, cmd.AccessToken)
	if err != nil {
		slog.Error("Failed to create code host client for profile fetch", "error", err)
		return nil, err
	}

	profile, err := ch.GetAuthenticatedUser(ctx)
	if err != nil {
		slog.Error("Failed to fetch authenticated user profile", "error", err)
		return nil, err
	}

	slog.Info("GetUserProfile query completed", "user_login", profile.Login, "user_id", profile.ID)
	return profile, nil
}

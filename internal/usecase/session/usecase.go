package session

import (
	"context"

	"2024_1_kayros/internal/repository/session"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error
	DeleteKey(ctx context.Context, key alias.SessionKey) error
}

type UsecaseLayer struct {
	repoSession session.Repo
}

func NewUsecaseLayer(repoSessionProps session.Repo) Usecase {
	return &UsecaseLayer{
		repoSession: repoSessionProps,
	}
}

func (uc *UsecaseLayer) GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error) {
	return uc.repoSession.GetValue(ctx, key)
}

func (uc *UsecaseLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error {
	return uc.repoSession.SetValue(ctx, key, value)
}

func (uc *UsecaseLayer) DeleteKey(ctx context.Context, key alias.SessionKey) error {
	return uc.repoSession.DeleteKey(ctx, key)
}

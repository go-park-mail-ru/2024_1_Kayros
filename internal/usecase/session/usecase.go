package session

import (
	"context"

	"2024_1_kayros/internal/repository/session"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetValue(context.Context, alias.SessionKey) (alias.SessionValue, error)
	SetValue(context.Context, alias.SessionKey, alias.SessionValue) (bool, error)
	DeleteKey(context.Context, alias.SessionKey) (bool, error)
}

type UsecaseLayer struct {
	repoUser    user.Repo
	repoSession session.Repo
}

func NewUsecaseLayer(repoUserProps user.Repo, repoSessionProps session.Repo) Usecase {
	return &UsecaseLayer{
		repoUser:    repoUserProps,
		repoSession: repoSessionProps,
	}
}

func (uc *UsecaseLayer) GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error) {
	value, err := uc.repoSession.GetValue(ctx, key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (uc *UsecaseLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) (bool, error) {
	wasSet, err := uc.repoSession.SetValue(ctx, key, value)
	return wasSet, err
}

func (uc *UsecaseLayer) DeleteKey(ctx context.Context, key alias.SessionKey) (bool, error) {
	wasDeleted, err := uc.repoSession.DeleteKey(ctx, key)
	return wasDeleted, err
}

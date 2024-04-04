package session

import (
	"context"

	"2024_1_kayros/internal/repository/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"go.uber.org/zap"
)

type Usecase interface {
	GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error
	DeleteKey(ctx context.Context, key alias.SessionKey) (bool, error)
}

type UsecaseLayer struct {
	repoSession session.Repo
	logger      *zap.Logger
}

func NewUsecaseLayer(repoSessionProps session.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoSession: repoSessionProps,
		logger:      loggerProps,
	}
}

func (uc *UsecaseLayer) GetValue(ctx context.Context, key alias.SessionKey) (alias.SessionValue, error) {
	methodName := cnst.NameMethodSetValue
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	value, err := uc.repoSession.GetValue(ctx, key)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return value, err
}

func (uc *UsecaseLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue) error {
	methodName := cnst.NameMethodSetValue
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	err := uc.repoSession.SetValue(ctx, key, value)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return err
}

func (uc *UsecaseLayer) DeleteKey(ctx context.Context, key alias.SessionKey) (bool, error) {
	methodName := cnst.NameMethodSetValue
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	wasDeleted, err := uc.repoSession.DeleteKey(ctx, key)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return wasDeleted, err
}

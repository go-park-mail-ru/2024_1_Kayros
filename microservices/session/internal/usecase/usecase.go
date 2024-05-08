package usecase

import (
	"2024_1_kayros/config"
	"2024_1_kayros/microservices/session/internal/repo"
	sessionv1 "2024_1_kayros/microservices/session/proto"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)


type Usecase interface {
	sessionv1.UnsafeSessionManagerServer
	SetSession(ctx context.Context, data *sessionv1.SetSessionData) (*emptypb.Empty, error)
	GetSession(ctx context.Context, data *sessionv1.GetSessionData) (*sessionv1.SessionValue, error)
	DeleteSession(ctx context.Context, data *sessionv1.DeleteSessionData) (*emptypb.Empty, error)
}


type Layer struct{
	sessionv1.UnsafeSessionManagerServer
	repoCsrf repo.Repo
	repoSession repo.Repo
	cfg *config.Redis
}

func NewLayer (redisCsrfProps repo.Repo, redisSessionProps repo.Repo, cfgProps *config.Redis) Usecase {
	return Layer {
		repoCsrf: redisCsrfProps,
		repoSession: redisSessionProps,
		cfg: cfgProps,
	}
}

func (uc Layer) SetSession(ctx context.Context, data *sessionv1.SetSessionData) (*emptypb.Empty, error) {
	setData := &sessionv1.SetSessionPair {
		Key: data.GetKey(),
		Value: data.GetValue(),
	}
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return nil, uc.repoCsrf.SetValue(ctx, setData)
	}
	return nil, uc.repoSession.SetValue(ctx, setData)
}

func (uc Layer) GetSession(ctx context.Context, data *sessionv1.GetSessionData) (*sessionv1.SessionValue, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return uc.repoCsrf.GetValue(ctx, data.GetKey())
	}
	return uc.repoSession.GetValue(ctx, data.GetKey())
}

func (uc Layer) DeleteSession(ctx context.Context, data *sessionv1.DeleteSessionData) (*emptypb.Empty, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return nil, uc.repoCsrf.DeleteValue(ctx, data.GetKey())
	}
	return nil, uc.repoSession.DeleteValue(ctx, data.GetKey())
}
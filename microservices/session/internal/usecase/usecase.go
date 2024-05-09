package usecase

import (
	"2024_1_kayros/config"
	"2024_1_kayros/microservices/session/internal/repo"
	"2024_1_kayros/gen/go/session"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)


type Usecase interface {
	session.UnsafeSessionManagerServer
	SetSession(ctx context.Context, data *session.SetSessionData) (*emptypb.Empty, error)
	GetSession(ctx context.Context, data *session.GetSessionData) (*session.SessionValue, error)
	DeleteSession(ctx context.Context, data *session.DeleteSessionData) (*emptypb.Empty, error)
}


type Layer struct{
	session.UnsafeSessionManagerServer
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

func (uc Layer) SetSession(ctx context.Context, data *session.SetSessionData) (*emptypb.Empty, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return nil, uc.repoCsrf.SetValue(ctx, data.GetKey(), data.GetValue())
	}
	return nil, uc.repoSession.SetValue(ctx, data.GetKey(), data.GetValue())
}

func (uc Layer) GetSession(ctx context.Context, data *session.GetSessionData) (*session.SessionValue, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return uc.repoCsrf.GetValue(ctx, data.GetKey())
	}
	return uc.repoSession.GetValue(ctx, data.GetKey())
}

func (uc Layer) DeleteSession(ctx context.Context, data *session.DeleteSessionData) (*emptypb.Empty, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		return nil, uc.repoCsrf.DeleteValue(ctx, data.GetKey())
	}
	return nil, uc.repoSession.DeleteValue(ctx, data.GetKey())
}
package usecase

import (
	"2024_1_kayros/config"
	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/microservices/session/internal/repo"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
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
		err := uc.repoCsrf.SetValue(ctx, data.GetKey(), data.GetValue())
		if err != nil {
			return nil, grpcerr.NewError(codes.Internal, err.Error())
		}
	}
	err := uc.repoSession.SetValue(ctx, data.GetKey(), data.GetValue())
	if err != nil {
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	return nil, nil
}

func (uc Layer) GetSession(ctx context.Context, data *session.GetSessionData) (*session.SessionValue, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		value, err := uc.repoCsrf.GetValue(ctx, data.GetKey())
		if err != nil {
			if errors.Is(err, myerrors.RedisNoData) {
				return &session.SessionValue{}, grpcerr.NewError(codes.NotFound, err.Error())
			}
			return &session.SessionValue{}, grpcerr.NewError(codes.Internal, err.Error())
		}
		return value, nil
	}
	value, err :=  uc.repoCsrf.GetValue(ctx, data.GetKey())
	if err != nil {
		if errors.Is(err, myerrors.RedisNoData) {
			return &session.SessionValue{}, grpcerr.NewError(codes.NotFound, err.Error())
		}
	 	return &session.SessionValue{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return value, nil
}

func (uc Layer) DeleteSession(ctx context.Context, data *session.DeleteSessionData) (*emptypb.Empty, error) {
	if data.GetDatabase() == int32(uc.cfg.DatabaseCsrf) {
		err := uc.repoCsrf.DeleteValue(ctx, data.GetKey())
		if err != nil {
			return nil, grpcerr.NewError(codes.Internal, err.Error())
		}
	}
	err := uc.repoSession.DeleteValue(ctx, data.GetKey())
	if err != nil {
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	return nil, nil
}
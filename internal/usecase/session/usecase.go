package session

import (
	"2024_1_kayros/gen/go/session"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
)

//go:generate mockgen -source=./usecase.go -destination=./usecase_mock.go -package=session
type Usecase interface {
	GetValue(ctx context.Context, key alias.SessionKey, databaseNum int32) (alias.SessionValue, error)
	SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, databaseNum int32) error
	DeleteKey(ctx context.Context, key alias.SessionKey, databaseNum int32) error
}

type UsecaseLayer struct {
	client session.SessionManagerClient
}

func NewUsecaseLayer(clientProps session.SessionManagerClient) Usecase {
	return &UsecaseLayer{
		client: clientProps,
	}
}

func (uc *UsecaseLayer) GetValue(ctx context.Context, key alias.SessionKey, databaseNum int32) (alias.SessionValue, error) {
	data := &session.GetSessionData{
		Key:      string(key),
		Database: databaseNum,
	}
	value, err := uc.client.GetSession(ctx, data)
	if err != nil {
		fmt.Printf("%v", err)
		if grpcerr.Is(err, codes.NotFound, myerrors.RedisNoData) {
			return "", myerrors.RedisNoData
		}
		return "", err
	}
	return alias.SessionValue(value.GetData()), err
}

func (uc *UsecaseLayer) SetValue(ctx context.Context, key alias.SessionKey, value alias.SessionValue, databaseNum int32) error {
	data := &session.SetSessionData{
		Key:      string(key),
		Value:    string(value),
		Database: databaseNum,
	}
	_, err := uc.client.SetSession(ctx, data)
	if err != nil {
		if grpcerr.Is(err, codes.Internal, myerrors.RedisNoData) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

func (uc *UsecaseLayer) DeleteKey(ctx context.Context, key alias.SessionKey, databaseNum int32) error {
	data := &session.DeleteSessionData{
		Key:      string(key),
		Database: databaseNum,
	}
	_, err := uc.client.DeleteSession(ctx, data)
	if err != nil {
		if grpcerr.Is(err, codes.Internal, myerrors.RedisNoData) {
			return myerrors.RedisNoData
		}
		return err
	}
	return nil
}

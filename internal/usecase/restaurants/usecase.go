package rest

//
//import (
//	"context"
//
//	"go.uber.org/zap"
//
//	"2024_1_kayros/internal/entity"
//	"2024_1_kayros/internal/repository/restaurants"
//	"2024_1_kayros/internal/utils/alias"
//	"2024_1_kayros/internal/utils/constants"
//	"2024_1_kayros/internal/utils/functions"
//)
//
//type Usecase interface {
//	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
//	GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error)
//}
//type UsecaseLayer struct {
//	repoRest restaurants.Repo
//	logger   *zap.Logger
//}
//
//func NewUsecaseLayer(repoRestProps restaurants.Repo, loggerProps *zap.Logger) Usecase {
//	return &UsecaseLayer{
//		repoRest: repoRestProps,
//		logger:   loggerProps,
//	}
//}
//
//func (uc *UsecaseLayer) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
//	methodName := constants.NameMethodGetAllRests
//	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
//	rests, err := uc.repoRest.GetAll(ctx, requestId)
//	if err != nil {
//		functions.LogUsecaseFail(uc.logger, requestId, methodName)
//		return nil, err
//	}
//	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
//	return rests, nil
//}
//
//func (uc *UsecaseLayer) GetById(ctx context.Context, restId alias.RestId) (*entity.Restaurant, error) {
//	methodName := constants.NameMethodGetRestById
//	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
//	rest, err := uc.repoRest.GetById(ctx, requestId, restId)
//	if err != nil {
//		functions.LogUsecaseFail(uc.logger, requestId, methodName)
//		return nil, err
//	}
//	functions.LogOk(uc.logger, requestId, methodName, constants.UsecaseLayer)
//	return rest, nil
//}

package usecase

import (
	"context"
	"sync"

	"2024_1_kayros/microservices/restaurants/internal/repo"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

type Rest interface {
	GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error)
	GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error)
	CreateComment(context.Context, *rest.Comment) (*rest.Comment, error)
	DeleteComment(context.Context, *rest.CommentId) (*rest.Empty, error)
	GetCommentsByRest(context.Context, *rest.RestId) (*rest.CommentList, error)
}

type RestLayer struct {
	rest.UnimplementedRestWorkerServer

	mu       *sync.RWMutex
	repoRest repo.Rest
}

func NewRestLayer(repoRestProps repo.Rest) *RestLayer {
	return &RestLayer{
		UnimplementedRestWorkerServer: rest.UnimplementedRestWorkerServer{},
		repoRest:                      repoRestProps,
		mu:                            &sync.RWMutex{},
	}
}

func (uc *RestLayer) GetAll(ctx context.Context, _ *rest.Empty) (*rest.RestList, error) {
	rests, err := uc.repoRest.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return rests, nil
}

func (uc *RestLayer) GetById(ctx context.Context, id *rest.RestId) (*rest.Rest, error) {
	rest, err := uc.repoRest.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return rest, nil
}

func (uc *RestLayer) CreateComment(context.Context, *rest.Comment) (*rest.Comment, error) {
	return nil, nil
}
func (uc *RestLayer) DeleteComment(context.Context, *rest.CommentId) (*rest.Empty, error) {
	return nil, nil
}
func (uc *RestLayer) GetCommentsByRest(context.Context, *rest.RestId) (*rest.CommentList, error) {
	return nil, nil
}

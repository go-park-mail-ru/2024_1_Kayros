package statistic

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/statistic"
)

type Usecase interface {
	Create(ctx context.Context, questionId uint64, rating uint32, user string) error
	Update(ctx context.Context, questionId uint64, rating uint32, user string) error
	GetQuestionInfo(ctx context.Context) ([]*entity.Question, error)
	GetStatistic(ctx context.Context) ([]*entity.Statistic, error)
}

type UsecaseLayer struct {
	repoStatistic statistic.Repo
}

func NewUsecaseLayer(repoStatisticProps statistic.Repo) Usecase {
	return &UsecaseLayer{
		repoStatistic: repoStatisticProps,
	}
}

func (uc *UsecaseLayer) Create(ctx context.Context, questionId uint64, rating uint32, user string) error {
	err := uc.repoStatistic.Create(ctx, questionId, rating, user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) Update(ctx context.Context, questionId uint64, rating uint32, user string) error {
	err := uc.repoStatistic.Update(ctx, questionId, rating, user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) GetQuestionInfo(ctx context.Context) ([]*entity.Question, error) {
	qs, err := uc.repoStatistic.GetQuestionInfo(ctx)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func (uc *UsecaseLayer) GetStatistic(ctx context.Context) ([]*entity.Statistic, error) {
	stats, err := uc.repoStatistic.GetStatistic(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

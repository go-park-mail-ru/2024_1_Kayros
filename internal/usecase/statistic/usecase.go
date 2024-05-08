package statistic

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/statistic"
)

type Usecase interface {
	Create(ctx context.Context, questionId uint64, rating uint8, userId string) error
	GetQuestionsOnFocus(ctx context.Context, url string) ([]*entity.Question, error)
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

func (uc *UsecaseLayer) Create(ctx context.Context, questionId uint64, rating uint8, userId string) error {
	err := uc.repoStatistic.Create(ctx, questionId, rating, userId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UsecaseLayer) GetQuestionsOnFocus(ctx context.Context, url string) ([]*entity.Question, error) {
	qs, err := uc.repoStatistic.GetQuestionsOnFocus(ctx, url)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func (uc *UsecaseLayer) GetStatistic(ctx context.Context) ([]*entity.Statistic, error) {
	qs, err := uc.repoStatistic.GetQuestions(ctx)
	if err != nil {
		return nil, err
	}
	stats := []*entity.Statistic{}
	for _, q := range qs {
		st := &entity.Statistic{QuestionId: q.Id, QuestionName: q.Name}
		if q.ParamType == "NPS" {
			res, err := uc.repoStatistic.NPS(ctx, q.Id)
			if err != nil {
				return nil, err
			}
			st.NPS = res
			stats = append(stats, st)
		}
		if q.ParamType == "CSAT" {
			res, err := uc.repoStatistic.CSAT(ctx, q.Id)
			if err != nil {
				return nil, err
			}
			st.CSAT = res
			stats = append(stats, st)
		}
	}
	return stats, nil
}

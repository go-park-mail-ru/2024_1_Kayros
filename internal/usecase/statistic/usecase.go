package statistic

import (
	"context"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/statistic"
)

type Usecase interface {
	Create(ctx context.Context, questionId uint64, rating uint32, user string) error
	Update(ctx context.Context, questionId uint64, rating uint32, user string) error
	GetQuestionInfo(ctx context.Context, url string) ([]*entity.Question, error)
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
	currentTime := time.Now().UTC()
	timeForDB := currentTime.Format("2006-01-02T15:04:05Z07:00")
	err := uc.repoStatistic.Create(ctx, questionId, rating, user, timeForDB)
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

func (uc *UsecaseLayer) GetQuestionInfo(ctx context.Context, url string) ([]*entity.Question, error) {
	qs, err := uc.repoStatistic.GetQuestionInfo(ctx, url)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func (uc *UsecaseLayer) GetStatistic(ctx context.Context) ([]*entity.Statistic, error) {
	//stats, err := uc.repoStatistic.GetStatistic(ctx)
	//if err != nil {
	//	return nil, err
	//}
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

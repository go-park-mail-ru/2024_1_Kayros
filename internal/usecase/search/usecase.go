package search

import (
	"context"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/search"
)

type Usecase interface {
	Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error)
}

type UsecaseLayer struct {
	repoSearch search.Repo
}

func NewUsecaseLayer(repoSearchProps search.Repo) Usecase {
	return &UsecaseLayer{
		repoSearch: repoSearchProps,
	}
}

func (uc *UsecaseLayer) Search(ctx context.Context, search string) ([]*dto.RestaurantAndFood, error) {
	rests, err := uc.repoSearch.Search(ctx, search)
	if err != nil {
		return nil, err
	}
	return rests, nil
}

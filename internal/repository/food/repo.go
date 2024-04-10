package food

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

const NoFoodError = "Такого блюда нет"

type Repo interface {
	GetByRestId(ctx context.Context, requestId string, restId alias.RestId) ([]*entity.Food, error)
	GetById(ctx context.Context, requestId string, foodId alias.FoodId) (*entity.Food, error)
}

type RepoLayer struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRepoLayer(dbProps *sql.DB, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		db:     dbProps,
		logger: loggerProps,
	}
}

func (repo *RepoLayer) GetByRestId(ctx context.Context, requestId string, restId alias.RestId) ([]*entity.Food, error) {
	rows, err := repo.db.QueryContext(ctx,
		`SELECT c.name, f.id, f.name, description, restaurant_id, weight, price, img_url FROM food as f 
    JOIN category as c ON f.category_id=c.id WHERE restaurant_id = $1 ORDER BY category_id`, uint64(restId))
	if err != nil {
		functions.LogError(repo.logger, requestId, constants.NameMethodGetFoodByRest, err, constants.RepoLayer)
		return nil, err
	}
	food := []*entity.Food{}
	for rows.Next() {
		item := entity.Food{}
		err = rows.Scan(&item.Category, &item.Id, &item.Name, &item.Description, &item.RestaurantId,
			&item.Weight, &item.Price, &item.ImgUrl)
		if err != nil {
			functions.LogError(repo.logger, requestId, constants.NameMethodGetFoodByRest, err, constants.RepoLayer)
			return nil, err
		}
		food = append(food, &item)
	}
	functions.LogOk(repo.logger, requestId, constants.NameMethodGetFoodByRest, constants.RepoLayer)
	return food, nil
}

func (repo *RepoLayer) GetById(ctx context.Context, requestId string, foodId alias.FoodId) (*entity.Food, error) {
	row := repo.db.QueryRowContext(ctx,
		`SELECT id, name, description, restaurant_id, category_id, weight, price, img_url 
				FROM food WHERE id=$1`, uint64(foodId))
	var item entity.Food
	err := row.Scan(&item.Id, &item.Name, &item.Description, &item.RestaurantId,
		&item.Category, &item.Weight, &item.Price, &item.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		functions.LogWarn(repo.logger, requestId, constants.NameMethodGetFoodByRest, fmt.Errorf(NoFoodError), constants.RepoLayer)
		return nil, fmt.Errorf(NoFoodError)
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, constants.NameMethodGetFoodById, err, constants.RepoLayer)
		return nil, err
	}
	functions.LogOk(repo.logger, requestId, constants.NameMethodGetFoodByRest, constants.RepoLayer)
	return &item, nil
}

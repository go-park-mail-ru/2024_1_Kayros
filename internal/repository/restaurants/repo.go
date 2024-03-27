package restaurants

import (
	"context"
	"database/sql"

	"2024_1_kayros/internal/entity"
)

type RestaurantRepoInterface interface {
	GetAll(ctx context.Context) ([]*entity.Restaurant, error)
	GetById(ctx context.Context, id int) (*entity.Restaurant, error)
}

type RestaurantRepo struct {
	DB *sql.DB
}

func NewRestaurantRepo(db *sql.DB) *RestaurantRepo {
	return &RestaurantRepo{DB: db}
}

func (repo *RestaurantRepo) GetAll(ctx context.Context) ([]*entity.Restaurant, error) {
	var rests []*entity.Restaurant
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, name, short_description, img_url FROM Restaurant")
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rest := &entity.Restaurant{}
		err = rows.Scan(&rest.Id, &rest.Name, &rest.ShortDescription, &rest.ImgUrl)
		if err != nil {
			return nil, err
		}
		rests = append(rests, rest)
	}
	return rests, nil
}

func (repo *RestaurantRepo) GetById(ctx context.Context, id int) (*entity.Restaurant, error) {
	rest := &entity.Restaurant{}
	row := repo.DB.QueryRowContext(ctx, "SELECT id, name, long_description, img_url FROM Restaurant WHERE id=$1", id)
	err := row.Scan(&rest.Id, &rest.Name, &rest.LongDescription, &rest.ImgUrl)
	if err != nil {
		return nil, err
	}
	return rest, nil
}

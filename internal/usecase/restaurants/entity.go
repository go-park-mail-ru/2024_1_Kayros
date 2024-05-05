package rest

import (
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	rest "2024_1_kayros/microservices/restaurants/proto"
)

func FromGrpcStructToRestaurant(grpcRest *rest.Rest) *entity.Restaurant {
	return &entity.Restaurant{
		Id:               grpcRest.Id,
		Name:             grpcRest.Name,
		ShortDescription: grpcRest.ShortDescription,
		LongDescription:  grpcRest.LongDescription,
		ImgUrl:           grpcRest.ImgUrl,
	}
}

func FromGrpcStructToRestaurantArray(grpcRest *rest.RestList) []*entity.Restaurant {
	if len(grpcRest.GetRest()) == 0 {
		return nil
	}
	restArray := make([]*entity.Restaurant, len(grpcRest.GetRest()))
	for i, r := range grpcRest.GetRest() {
		restArray[i] = FromGrpcStructToRestaurant(r)
	}
	return restArray
}

func FromGrpcStructToCategory(grpcCat *rest.Category) *entity.Category {
	return &entity.Category{
		Id:   alias.CategoryId(grpcCat.Id),
		Name: grpcCat.Name,
	}
}

func FromGrpcStructToCategoryArray(grpcCat *rest.CategoryList) []*entity.Category {
	if len(grpcCat.GetC()) == 0 {
		return nil
	}
	categoryArray := make([]*entity.Category, len(grpcCat.GetC()))
	for i, r := range grpcCat.GetC() {
		categoryArray[i] = FromGrpcStructToCategory(r)
	}
	return categoryArray
}
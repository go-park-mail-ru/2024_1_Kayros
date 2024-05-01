package rest

import (
	"2024_1_kayros/internal/entity"
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

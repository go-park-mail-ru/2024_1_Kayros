package entity

import rest "2024_1_kayros/microservices/restaurants/proto"

type Restaurant struct {
	Id               uint64
	Name             string
	ShortDescription string
	LongDescription  string
	Address          string
	ImgUrl           string
}

func FromGrpcStructToRestaurant(grpcRest *rest.Rest) *Restaurant {
	return &Restaurant{
		Id:               grpcRest.Id,
		Name:             grpcRest.Name,
		ShortDescription: grpcRest.ShortDescription,
		LongDescription:  grpcRest.LongDescription,
		ImgUrl:           grpcRest.ImgUrl,
	}
}

func FromGrpcStructToRestaurantArray(grpcRest *rest.RestList) []*Restaurant {
	if len(grpcRest.GetRest()) == 0 {
		return nil
	}
	restArray := make([]*Restaurant, len(grpcRest.GetRest()))
	for i, r := range grpcRest.GetRest() {
		restArray[i] = FromGrpcStructToRestaurant(r)
	}
	return restArray
}

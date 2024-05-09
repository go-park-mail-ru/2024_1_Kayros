package entity

import (
	"2024_1_kayros/gen/go/food"
	"2024_1_kayros/internal/utils/alias"
)

type Category struct {
	Id   alias.CategoryId
	Name string
	Food []*Food
}

type Food struct {
	Id           uint64
	Name         string
	Description  string
	RestaurantId uint64
	Category     string
	Weight       uint64
	Price        uint64
	ImgUrl       string
}

type FoodInOrder struct {
	Id           uint64
	Name         string
	Weight       uint64
	Price        uint64
	Count        uint64
	ImgUrl       string
	RestaurantId uint64
}

func CnvProtoFoodIntoEntityFood(food *food.Food) *Food {
	return &Food {
		Id: food.GetId(),
		Name: food.GetName(),
		Description: food.GetDescription(),
		RestaurantId: food.GetRestaurantId(),
		Category: food.GetCategory(),
		Weight: food.GetWeight(),
		Price: food.GetPrice(),
		ImgUrl: food.GetImgUrl(),
	}
}

func CnvEntityFoodIntoProtoFood(f *Food) *food.Food {
	return &food.Food {
		Id: f.Id,
		Name: f.Name,
		Description: f.Description,
		RestaurantId: f.RestaurantId,
		Category: f.Category,
		Weight: f.Weight,
		Price: f.Price,
		ImgUrl: f.ImgUrl,
	}
}

func CnvEntityCtgIntoProtoCtgs(ctgs []*Category) *food.RestCategories {
	protoCtgs := []*food.Category{}
	for _, ctg := range ctgs {
		protoCtg := food.Category{}
		protoFood := []*food.Food{}
		for _, f := range ctg.Food {
			protoFood = append(protoFood, CnvEntityFoodIntoProtoFood(f))
		}
		protoCtg.Food = protoFood
		protoCtg.Id = uint64(ctg.Id)
		protoCtg.Name = ctg.Name
		protoCtgs = append(protoCtgs, &protoCtg)
	}
	return &food.RestCategories{Category: protoCtgs}
}

func CnvProtoCtgIntoEntityCtg(protoCtgs *food.RestCategories) []*Category {
	ctgs := []*Category{}
	for _, protoCtg := range protoCtgs.Category {
		ctg := Category{}
		food := []*Food{}
		for _, f := range protoCtg.Food {
			food = append(food, CnvProtoFoodIntoEntityFood(f))
		}
		ctg.Food = food
		ctg.Id = alias.CategoryId(protoCtg.Id)
		ctg.Name = protoCtg.Name
		ctgs = append(ctgs, &ctg)
	}
	return ctgs
}
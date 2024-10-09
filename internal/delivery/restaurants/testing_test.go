package restaurants

import (
	"testing"

	food "2024_1_kayros/internal/usecase/food/mocks"
	rest "2024_1_kayros/internal/usecase/restaurants/mocks"
	"2024_1_kayros/internal/usecase/user"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type testFixtures struct {
	ctrl       *gomock.Controller
	handler    *RestaurantHandler
	mockUcRest *rest.MockUsecase
	mockUcFood *food.MockUsecase
}

func setUp(t *testing.T) testFixtures {
	ctrl := gomock.NewController(t)
	mockUcRest := rest.NewMockUsecase(ctrl)
	mockUcFood := food.NewMockUsecase(ctrl)
	logger := zap.NewNop()

	return testFixtures{
		ctrl:       ctrl,
		handler:    NewRestaurantHandler(mockUcRest, mockUcFood, &user.UsecaseLayer{}, logger),
		mockUcRest: mockUcRest,
		mockUcFood: mockUcFood,
	}
}

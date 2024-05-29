package usecase

import (
	"testing"

	repo "2024_1_kayros/microservices/restaurants/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type testFixtures struct {
	ctrl     *gomock.Controller
	layer    *RestLayer
	mockRepo *repo.MockRest
}

func setUp(t *testing.T) testFixtures {
	ctrl := gomock.NewController(t)
	mockRepo := repo.NewMockRest(ctrl)

	logger := zap.NewNop()

	return testFixtures{
		ctrl:     ctrl,
		layer:    NewRestLayer(mockRepo, logger),
		mockRepo: mockRepo,
	}
}

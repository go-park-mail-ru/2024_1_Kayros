package usecase

import (
	"testing"

	user "2024_1_kayros/gen/go/user/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type testFixtures struct {
	ctrl           *gomock.Controller
	layer          Usecase
	mockUserClient *user.MockUserManagerClient
}

func setUp(t *testing.T) testFixtures {
	ctrl := gomock.NewController(t)
	mockUserClient := user.NewMockUserManagerClient(ctrl)

	logger := zap.NewNop()

	return testFixtures{
		ctrl:           ctrl,
		layer:          NewLayer(mockUserClient, logger),
		mockUserClient: mockUserClient,
	}
}

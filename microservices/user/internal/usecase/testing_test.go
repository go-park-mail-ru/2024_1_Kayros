package usecase

import (
	"testing"

	minio "2024_1_kayros/internal/repository/minios3/mocks"
	user "2024_1_kayros/microservices/user/internal/repo/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type testFixtures struct {
	ctrl      *gomock.Controller
	layer     Usecase
	userRepo  *user.MockRepo
	minioRepo *minio.MockRepo
}

func setUp(t *testing.T) testFixtures {
	ctrl := gomock.NewController(t)
	minioRepo := minio.NewMockRepo(ctrl)
	userRepo := user.NewMockRepo(ctrl)
	logger := zap.NewNop()

	return testFixtures{
		ctrl:      ctrl,
		layer:     NewLayer(userRepo, minioRepo, logger),
		userRepo:  userRepo,
		minioRepo: minioRepo,
	}
}

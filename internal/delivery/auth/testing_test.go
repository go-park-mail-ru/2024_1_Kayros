package auth

import (
	"io"
	"testing"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	auth "2024_1_kayros/internal/usecase/auth/mocks"
	session "2024_1_kayros/internal/usecase/session/mocks"
	"2024_1_kayros/internal/utils/functions"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

type testFixtures struct {
	ctrl       *gomock.Controller
	handler    *Delivery
	mockUcAuth *auth.MockUsecase
	mockUcSess *session.MockUsecase
}

type errorReader struct{}

func (er *errorReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func setUp(t *testing.T) testFixtures {
	ctrl := gomock.NewController(t)
	mockUcAuth := auth.NewMockUsecase(ctrl)
	mockUcSess := session.NewMockUsecase(ctrl)
	logger := zap.NewNop()
	functions.InitDtoValidator(logger)
	return testFixtures{
		ctrl:       ctrl,
		handler:    NewDeliveryLayer(&config.Project{}, mockUcSess, mockUcAuth, logger, &metrics.Metrics{}),
		mockUcAuth: mockUcAuth,
		mockUcSess: mockUcSess,
	}
}

package session

import (
	"context"
	"testing"
	"time"

	"2024_1_kayros/internal/utils/alias"
	"github.com/golang/mock/gomock"
)

const place = "session.usecase_test."

func TestUser_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)

	ctx, errCtx := context.WithTimeout(context.Background(), 5*time.Second)
	if errCtx != nil {
		t.Errorf("%sGetValue err: %v", place, errCtx)
	}

	sessionKey := alias.SessionKey("c05c5ad9-1bf8-436d-89aa-ee8321d3a41a")
	mockSession := NewMockUsecase(ctl)
	gomock.InOrder(
		mockSession.EXPECT().GetValue(ctx, sessionKey).Return(nil),
	)

	//logger := zap.Must(zap.NewProduction())
	//ucSession := UsecaseLayer{mockSession, logger}
	//sessionValue, err := ucSession.GetValue(ctx, sessionKey)
	//if err != nil {
	//	t.Errorf("%sGetValue err: %v", place, err)
	//}

}

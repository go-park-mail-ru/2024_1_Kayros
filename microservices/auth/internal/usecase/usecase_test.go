package usecase

import (
	"context"
	"fmt"
	"testing"

	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("error in checking if user exists", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUserClient.EXPECT().GetData(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))
		_, err := s.layer.SignUp(ctx, &auth.SignUpCredentials{})
		assert.Error(t, err)
	})

	t.Run("user already exists", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUserClient.EXPECT().GetData(ctx, gomock.Any()).Return(&user.User{}, nil)
		_, err := s.layer.SignUp(ctx, &auth.SignUpCredentials{})
		assert.Error(t, err)
	})



	t.Run("user creation error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUserClient.EXPECT().GetData(ctx, gomock.Any()).Return(nil, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsUserRelation.Error()))
		s.mockUserClient.EXPECT().Create(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))

		_, err := s.layer.SignUp(ctx, &auth.SignUpCredentials{})
		assert.Error(t, err)
	 })

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()

		s.mockUserClient.EXPECT().GetData(ctx, gomock.Any()).Return(nil, grpcerr.NewError(codes.NotFound, myerrors.SqlNoRowsUserRelation.Error()))
		s.mockUserClient.EXPECT().Create(ctx, gomock.Any()).Return(&user.User{Id: 1}, nil)

		u, err := s.layer.SignUp(ctx, &auth.SignUpCredentials{})
		assert.NoError(t, err)
		assert.Equal(t, &auth.User{Id: 1}, u)
	})

}

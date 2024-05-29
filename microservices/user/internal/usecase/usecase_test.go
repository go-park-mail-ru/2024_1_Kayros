package usecase

import (
	"context"
	"fmt"
	"testing"

	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestIsPasswordEquals(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		} ..

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(nil, fmt.Errorf("error"))
		bVal, err := s.layer.IsPassswordEquals(ctx, passCheck)
		assert.Equal(t, &wrapperspb.BoolValue{Value: false}, bVal)
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(nil, myerrors.SqlNoRowsUserRelation)
		bVal, err := s.layer.IsPassswordEquals(ctx, passCheck)
		assert.Equal(t, &wrapperspb.BoolValue{Value: false}, bVal)
		assert.Equal(t, &grpcerr.Error{Status: codes.NotFound, Message: myerrors.SqlNoRowsUserRelation.Error()}, err)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(&user.User{}, nil)
		_, err := s.layer.IsPassswordEquals(ctx, passCheck)
		assert.NoError(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(nil, fmt.Errorf("error"))
		u, err := s.layer.Create(ctx, &user.User{Email: passCheck.Email})
		assert.Equal(t, &user.User{}, u)
		assert.Error(t, err)
	})

	t.Run("already exists", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(&user.User{}, nil)
		u, err := s.layer.Create(ctx, &user.User{Email: passCheck.Email})
		assert.Equal(t, &user.User{}, u)
		assert.Error(t, err)
	})

	t.Run("error in creating", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}
		u := &user.User{Email: passCheck.Email}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(nil, myerrors.SqlNoRowsUserRelation)
		s.userRepo.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("error"))
		userCreated, err := s.layer.Create(ctx, u)
		assert.Error(t, err)
		assert.Equal(t, &user.User{}, userCreated)
	})
}

func TestGetData(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(nil, fmt.Errorf("error"))
		u, err := s.layer.GetData(ctx, &user.Email{Email: passCheck.Email})
		assert.Equal(t, &user.User{}, u)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		passCheck := &user.PasswordCheck{
			Email:    "aaa@aa.aa",
			Password: "aaaaaaaa",
		}

		s.userRepo.EXPECT().GetByEmail(ctx, &user.Email{Email: passCheck.Email}).Return(&user.User{}, nil)
		u, err := s.layer.GetData(ctx, &user.Email{Email: passCheck.Email})
		assert.Equal(t, &user.User{}, u)
		assert.NoError(t, err)
	})

}

func TestGetAddressByUnauthId(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("internal error", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		uId := &user.UnauthId{UnauthId: "a"}

		s.userRepo.EXPECT().GetAddressByUnauthId(ctx, uId).Return(nil, fmt.Errorf("error"))
		u, err := s.layer.GetAddressByUnauthId(ctx, uId)
		assert.Equal(t, &user.Address{}, u)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		s := setUp(t)
		defer s.ctrl.Finish()
		uId := &user.UnauthId{UnauthId: "a"}
		adr := &user.Address{Address: "aaa"}

		s.userRepo.EXPECT().GetAddressByUnauthId(ctx, uId).Return(adr, nil)
		u, err := s.layer.GetAddressByUnauthId(ctx, uId)
		assert.Equal(t, adr, u)
		assert.NoError(t, err)
	})
}

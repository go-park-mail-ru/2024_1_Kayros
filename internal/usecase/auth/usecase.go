package auth

import (
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"context"

	"google.golang.org/grpc/codes"
)

type Usecase interface {
	SignUp(ctx context.Context, u *entity.User, unauthId string) (*entity.User, error)
	SignIn(ctx context.Context, email string, password string, unauthId string) (*entity.User, error)
}

type UsecaseLayer struct {
	grpcClient auth.AuthManagerClient
}

func NewUsecaseLayer(restClientProps auth.AuthManagerClient) Usecase {
	return &UsecaseLayer{
		grpcClient: restClientProps,
	}
}

func (uc *UsecaseLayer) SignUp(ctx context.Context, u *entity.User, unauthId string) (*entity.User, error) {
	data := &auth.SignUpCredentials{
		Email:    u.Email,
		UnauthId: unauthId,
		Password: u.Password,
		Name:     u.Name,
	}
	uSignedUp, err := uc.grpcClient.SignUp(ctx, data)
	if err != nil {
		if grpcerr.Is(err, codes.AlreadyExists, myerrors.UserAlreadyExist) {
			return &entity.User{}, myerrors.UserAlreadyExist
		}
		if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUserRelation) {
			return &entity.User{}, myerrors.SqlNoRowsUserRelation
		}
	}
	return cnvAuthUserIntoEntityUser(uSignedUp), nil
}

func cnvAuthUserIntoEntityUser(u *auth.User) *entity.User {
	return &entity.User{
		Id:         u.GetId(),
		Name:       u.GetName(),
		Phone:      u.GetPhone(),
		Email:      u.GetEmail(),
		Address:    u.GetAddress(),
		ImgUrl:     u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
		Password:   u.GetPassword(),
	}
}

func (uc *UsecaseLayer) SignIn(ctx context.Context, email string, password string, unauthId string) (*entity.User, error) {
	data := &auth.SignInCredentials{
		Email:    email,
		Password: password,
		UnauthId: unauthId,
	}
	u, err := uc.grpcClient.SignIn(ctx, data)
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return &entity.User{}, myerrors.SqlNoRowsUserRelation
		}
		if grpcerr.Is(err, codes.InvalidArgument, myerrors.BadAuthPassword) {
			return &entity.User{}, myerrors.BadAuthPassword
		}
		return &entity.User{}, err
	}
	return cnvAuthUserIntoEntityUser(u), nil
}

package auth

import (
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/internal/entity"
	"context"
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
	data := &auth.SignUpCredentials {
		Email: u.Email,
		UnauthId: unauthId,
		Password: u.Password,
		Name: u.Name,
	}
	uSignedUp, err := uc.grpcClient.SignUp(ctx, data)
	return cnvAuthUserIntoEntityUser(uSignedUp), err
}

func cnvAuthUserIntoEntityUser (u *auth.User) *entity.User {
	return &entity.User{
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email: u.GetEmail(),
		Address: u.GetAddress(),
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
		Password: u.GetPassword(),
	}
}

func (uc *UsecaseLayer) SignIn(ctx context.Context, email string, password string, unauthId string) (*entity.User, error) {
	data := &auth.SignInCredentials{
		Email: email,
		Password: password,
		UnauthId: unauthId,
	}
	u, err := uc.grpcClient.SignIn(ctx, data)
	return cnvAuthUserIntoEntityUser(u), err
}

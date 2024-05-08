package auth

import (
	"context"
	authv1 "2024_1_kayros/microservices/auth/proto"
)

type Usecase interface {
	SignUp(ctx context.Context, data *authv1.SignUpCredentials) (*authv1.User, error)
	SignIn(ctx context.Context, data *authv1.SignInCredentials) (*authv1.User, error)
}

type UsecaseLayer struct {
	grpcClient authv1.AuthManagerClient
}

func NewUsecaseLayer(restClientProps authv1.AuthManagerClient) Usecase {
	return &UsecaseLayer{
		grpcClient: restClientProps,
	}
}

func (uc *UsecaseLayer) SignUp(ctx context.Context, data *authv1.SignUpCredentials) (*authv1.User, error) {
	return uc.grpcClient.SignUp(ctx, data)
}

func (uc *UsecaseLayer) SignIn(ctx context.Context, data *authv1.SignInCredentials) (*authv1.User, error) {
	return uc.grpcClient.SignIn(ctx, data)
}

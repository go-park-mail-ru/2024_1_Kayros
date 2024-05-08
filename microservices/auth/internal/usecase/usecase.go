package usecase

import (
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	authv1 "2024_1_kayros/microservices/auth/proto"
	userv1 "2024_1_kayros/microservices/user/proto"
	"bytes"
	"context"
	"errors"
)


type Usecase interface {
	authv1.UnsafeAuthManagerServer
	SignUp(ctx context.Context, data *authv1.SignUpCredentials) (*authv1.User, error)
	SignIn(ctx context.Context, data *authv1.SignInCredentials) (*authv1.User, error)
}

type Layer struct {
	authv1.UnsafeAuthManagerServer
	client userv1.UserManagerClient
}

func NewLayer(clientProps userv1.UserManagerClient) Usecase {
	return &Layer{
		client: clientProps,
	}
}

func (uc *Layer) SignUp(ctx context.Context, data *authv1.SignUpCredentials) (*authv1.User, error) {
	isExist, err := uc.IsExistByEmail(ctx, &userv1.Email{Email: data.GetEmail()})
	if err != nil {
		// we can skip error `myerrors.SqlNoRowsUserRelation`, because user must not to be
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if isExist {
		return nil, myerrors.UserAlreadyExist
	}

	address, err := uc.client.GetAddressByUnauthId(ctx, &userv1.UnauthId{UnauthId: data.GetUnauthId()})
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
		return nil, err
	}
	data.SignUpData.Address = address.GetAddress()

	// we do copy for clean function
	uCopy := entity.Copy(convAuthUserIntoUser(data.GetSignUpData()))
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetSignUpData().GetPassword())
	uCopy.Password = &userv1.Password{Password: string(hashPassword)}

	uCreated, err := uc.client.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}
	return convUserIntoAuthUser(uCreated), nil
}

func (uc *Layer) SignIn(ctx context.Context, data *authv1.SignInCredentials) (*authv1.User, error) {
	u, err := uc.client.GetData(ctx, &userv1.Email{Email: data.GetEmail()})
	if err != nil {
		return nil, err
	}

	isEqual, err := uc.checkPassword(ctx, &userv1.Email{Email: data.Email}, &userv1.Password{Password: data.Password})
	if err != nil {
		return nil, err
	}
	if !isEqual {
		return nil, myerrors.BadAuthPassword
	}

	return convUserIntoAuthUser(u), nil
}

func (uc *Layer) IsExistByEmail(ctx context.Context, email *userv1.Email) (bool, error) {
	_, err := uc.client.GetData(ctx, email)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// checkPassword - method used to check password with password saved in database
func (uc *Layer) checkPassword(ctx context.Context, email *userv1.Email, password *userv1.Password) (bool, error) {
	u, err := uc.client.GetData(ctx, email)
	if err != nil {
		return false, err
	}
	uPasswordBytes := []byte(u.GetPassword().GetPassword())

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, password.GetPassword())
	return bytes.Equal(uPasswordBytes, hashPassword), nil
}

func convAuthUserIntoUser (u *authv1.User) *userv1.User {
	return &userv1.User {
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email: &userv1.Email{Email: u.GetEmail()},
		Address: &userv1.Address{Address: u.GetAddress()},
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
	}
}

func convUserIntoAuthUser (u *userv1.User) *authv1.User {
	return &authv1.User {
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email:  u.GetEmail().GetEmail(),
		Address: u.GetAddress().GetAddress(),
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
	}
}
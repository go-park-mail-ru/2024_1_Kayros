package usecase

import (
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/user"
	"context"
	"errors"
)


type Usecase interface {
	auth.UnsafeAuthManagerServer
	SignUp(ctx context.Context, data *auth.SignUpCredentials) (*auth.User, error)
	SignIn(ctx context.Context, data *auth.SignInCredentials) (*auth.User, error)
}

type Layer struct {
	auth.UnsafeAuthManagerServer
	client user.UserManagerClient
}

func NewLayer(clientProps user.UserManagerClient) Usecase {
	return &Layer{
		client: clientProps,
	}
}

func (uc *Layer) SignUp(ctx context.Context, data *auth.SignUpCredentials) (*auth.User, error) {
	isExist, err := uc.isExistByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		// we can skip error `myerrors.SqlNoRowsUserRelation`, because user must not to be
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if isExist {
		return nil, myerrors.UserAlreadyExist
	}

	address, err := uc.client.GetAddressByUnauthId(ctx, &user.UnauthId{UnauthId: data.GetUnauthId()})
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
		return nil, err
	}
	data.Address = address.GetAddress()

	// we do copy for clean function
	uCopy := entity.Copy(convAuthUserIntoUser(data))
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetPassword())
	uCopy.Password = string(hashPassword)

	uCreated, err := uc.client.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}
	return convUserIntoAuthUser(uCreated), nil
}

func (uc *Layer) SignIn(ctx context.Context, data *auth.SignInCredentials) (*auth.User, error) {
	u, err := uc.client.GetData(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		return nil, err
	}

	isEqual, err := uc.checkPassword(ctx, data.GetEmail(), data.GetPassword())
	if err != nil {
		return nil, err
	}
	if !isEqual {
		return nil, myerrors.BadAuthPassword
	}

	return convUserIntoAuthUser(u), nil
}

func (uc *Layer) isExistByEmail(ctx context.Context, email *user.Email) (bool, error) {
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
func (uc *Layer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
	passwordData := &user.PasswordCheck {
		Email: email,
		Password: password,
	}
	isEqual, err := uc.client.IsPassswordEquals(ctx, passwordData)
	return isEqual.Value, err
}

func convAuthUserIntoUser (u *auth.SignUpCredentials) *user.User {
	return &user.User {
		Name: u.GetName(),
		Email: u.GetEmail(),
		Address: u.GetAddress(),
		Password:  u.Password,
	}
}

func convUserIntoAuthUser (u *user.User) *auth.User {
	return &auth.User {
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email:  u.GetEmail(),
		Address: u.GetAddress(),
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
	}
}
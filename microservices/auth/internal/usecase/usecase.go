package usecase

import (
	"2024_1_kayros/gen/go/auth"
	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type Usecase interface {
	auth.UnsafeAuthManagerServer
	SignUp(ctx context.Context, data *auth.SignUpCredentials) (*auth.User, error)
	SignIn(ctx context.Context, data *auth.SignInCredentials) (*auth.User, error)
}

type Layer struct {
	auth.UnsafeAuthManagerServer
	client user.UserManagerClient
	logger *zap.Logger
}

func NewLayer(clientProps user.UserManagerClient, loggerProps *zap.Logger) Usecase {
	return &Layer{
		client: clientProps,
		logger: loggerProps,
	}
}

func (uc *Layer) SignUp(ctx context.Context, data *auth.SignUpCredentials) (*auth.User, error) {
	isExist, err := uc.isExistByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		return &auth.User{}, err
	}

	if isExist {
		uc.logger.Error(myerrors.UserAlreadyExist.Error())
		return &auth.User{}, grpcerr.NewError(codes.AlreadyExists, myerrors.UserAlreadyExist.Error())
	}

	address, err := uc.client.GetAddressByUnauthId(ctx, &user.UnauthId{UnauthId: data.GetUnauthId()})
	if err != nil && !grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUnauthAddressRelation) {
		uc.logger.Error(err.Error())
		return &auth.User{}, err
	}

	// we do copy for clean function
	uCopy := entity.Copy(convAuthUserIntoUser(data, address.GetAddress()))
	uCreated, err := uc.client.Create(ctx, uCopy)
	if err != nil {
		uc.logger.Error(err.Error())
		return &auth.User{}, err
	}

	return convUserIntoAuthUser(uCreated), nil
}

func (uc *Layer) SignIn(ctx context.Context, data *auth.SignInCredentials) (*auth.User, error) {
	u, err := uc.client.GetData(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		return &auth.User{}, err
	}

	isEqual, err := uc.checkPassword(ctx, data.GetEmail(), data.GetPassword())
	if err != nil {
		uc.logger.Error(err.Error())
		return &auth.User{}, err
	}
	if !isEqual {
		uc.logger.Error(myerrors.BadAuthPassword.Error())
		return &auth.User{}, grpcerr.NewError(codes.InvalidArgument, myerrors.BadAuthPassword.Error())
	}

	return convUserIntoAuthUser(u), nil
}

func (uc *Layer) isExistByEmail(ctx context.Context, email *user.Email) (bool, error) {
	_, err := uc.client.GetData(ctx, email)
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// checkPassword - method used to check password with password saved in database
func (uc *Layer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
	passwordData := &user.PasswordCheck{
		Email:    email,
		Password: password,
	}
	isEqual, err := uc.client.IsPassswordEquals(ctx, passwordData)
	return isEqual.Value, err
}

func convAuthUserIntoUser(u *auth.SignUpCredentials, address string) *user.User {
	return &user.User{
		Name:     u.GetName(),
		Email:    u.GetEmail(),
		Address:  address,
		Password: u.Password,
	}
}

func convUserIntoAuthUser(u *user.User) *auth.User {
	return &auth.User{
		Id:         u.GetId(),
		Name:       u.GetName(),
		Phone:      u.GetPhone(),
		Email:      u.GetEmail(),
		Address:    u.GetAddress(),
		ImgUrl:     u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
	}
}

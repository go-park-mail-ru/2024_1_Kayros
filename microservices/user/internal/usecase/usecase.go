package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/minios3"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/microservices/user/internal/repo"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/satori/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecase interface {
	user.UnsafeUserManagerServer
	GetData(ctx context.Context, email *user.Email) (*user.User, error)
	UpdateData(ctx context.Context, data *user.UpdateUserData) (*user.User, error)
	Create(ctx context.Context, data *user.User) (*user.User, error)
	UpdateAddress(ctx context.Context, data *user.AddressData) (*emptypb.Empty, error)
	SetNewPassword(ctx context.Context, data *user.PasswordsChange) (*emptypb.Empty, error)
	IsPassswordEquals(ctx context.Context, pwds *user.PasswordCheck) (*wrapperspb.BoolValue, error)
	UpdateAddressByUnauthId(ctx context.Context, unauth *user.AddressDataUnauth) (*emptypb.Empty, error)
	GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error)
}

type Layer struct {
	user.UnsafeUserManagerServer
	repoUser repo.Repo
	minio    minios3.Repo
	logger   *zap.Logger
}

func NewLayer(repoUserProps repo.Repo, repoMinio minios3.Repo, loggerProps *zap.Logger) Usecase {
	return &Layer{
		repoUser: repoUserProps,
		minio:    repoMinio,
		logger:   loggerProps,
	}
}

func deleteCredentials(u *user.User) *user.User {
	if u != nil {
		u.Password = ""
		u.CardNumber = ""
	}
	return u
}

// GetData - method calls repo method to receive user data.
func (uc *Layer) GetData(ctx context.Context, email *user.Email) (*user.User, error) {
	returnUser, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return &user.User{}, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return deleteCredentials(returnUser), nil
}

// UpdateData - method used to update user info. it accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc *Layer) UpdateData(ctx context.Context, data *user.UpdateUserData) (*user.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return &user.User{}, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	fillUserFields(u, data.GetUpdateInfo())
	if len(data.GetFileData()) != 0 && data.GetFileSize() != 0 {
		mimeType := functions.GetFileMimeType(data.GetFileData())
		if _, ok := cnst.ValidMimeTypes[mimeType]; !ok {
			uc.logger.Error(myerrors.WrongFileExtension.Error())
			return &user.User{}, grpcerr.NewError(codes.InvalidArgument, myerrors.WrongFileExtension.Error())
		}

		fileExtension := functions.GetFileExtension(data.GetFileName())
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, bytes.NewBuffer(data.GetFileData()), filename, data.GetFileSize(), mimeType)
		if err != nil {
			uc.logger.Error(err.Error())
			return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
		}
		u.ImgUrl = fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	}

	uc.logger.Info(fmt.Sprintf("%v", u))

	newEmail := &user.Email{Email: data.GetUpdateInfo().GetEmail()}
	uDB, err := uc.repoUser.GetByEmail(ctx, newEmail)
	fmt.Println(uDB == nil)
	uc.logger.Info(fmt.Sprintf("%v", uDB))
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	if uDB != nil && data.GetEmail() != newEmail.GetEmail() {
		uc.logger.Error(myerrors.UserAlreadyExist.Error())
		return &user.User{}, grpcerr.NewError(codes.AlreadyExists, myerrors.UserAlreadyExist.Error())
	}

	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return deleteCredentials(u), nil
}

// UpdateAddress - method updates only address.
func (uc *Layer) UpdateAddress(ctx context.Context, data *user.AddressData) (*emptypb.Empty, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}

	u.Address = data.GetAddress()
	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	return nil, nil
}

func (uc *Layer) GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error) {
	address, err := uc.repoUser.GetAddressByUnauthId(ctx, id)
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
			return &user.Address{}, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return &user.Address{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	return address, nil
}

// SetNewPassword - method used to set new password.
func (uc *Layer) SetNewPassword(ctx context.Context, data *user.PasswordsChange) (*emptypb.Empty, error) {
	// compare old password with database password
	isEqual, err := uc.checkPassword(ctx, data.GetEmail(), data.GetPassword())
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}

	if !isEqual {
		uc.logger.Error(myerrors.IncorrectCurrentPassword.Error())
		return nil, grpcerr.NewError(codes.InvalidArgument, myerrors.IncorrectCurrentPassword.Error())
	}
	if data.GetPassword() == data.GetNewPassword() {
		uc.logger.Error(myerrors.NewPassword.Error())
		return nil, grpcerr.NewError(codes.InvalidArgument, myerrors.NewPassword.Error())
	}

	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	hashPassword := functions.HashData(salt, data.GetNewPassword())
	u.Password = string(hashPassword)

	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		uc.logger.Error(err.Error())
		return nil, grpcerr.NewError(codes.Internal, err.Error())
	}
	return nil, nil
}

func (uc *Layer) UpdateAddressByUnauthId(ctx context.Context, unauth *user.AddressDataUnauth) (*emptypb.Empty, error) {
	err := uc.repoUser.UpdateAddressByUnauthId(ctx, unauth)
	if err != nil {
		uc.logger.Error(err.Error())
		if !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
			return nil, grpcerr.NewError(codes.Internal, err.Error())
		}
		err = uc.repoUser.CreateAddressByUnauthId(ctx, unauth)
		if err != nil {
			return nil, grpcerr.NewError(codes.Internal, err.Error())
		}
	}
	return nil, nil
}

// Create - method used to create new user in database.
func (uc *Layer) Create(ctx context.Context, data *user.User) (*user.User, error) {
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	hashPassword := functions.HashData(salt, data.GetPassword())

	uCopy := entity.Copy(data)
	uCopy.Password = string(hashPassword)

	_, err = uc.repoUser.GetByEmail(ctx, &user.Email{Email: uCopy.GetEmail()})
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}
	if err == nil {
		uc.logger.Error(myerrors.UserAlreadyExist.Error())
		return &user.User{}, grpcerr.NewError(codes.AlreadyExists, myerrors.UserAlreadyExist.Error())
	}
	if uCopy.ImgUrl == "" {
		uCopy.ImgUrl = "/minio-api/users/default.jpg"
	}
	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}

	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: uCopy.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		return &user.User{}, grpcerr.NewError(codes.Internal, err.Error())
	}

	return deleteCredentials(u), nil
}

// checkPassword - method used to check password with password saved in database
func (uc *Layer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: email})
	if err != nil {
		return false, err
	}
	uPasswordBytes := []byte(u.GetPassword())

	hashPassword := functions.HashData(uPasswordBytes[0:8], password)
	return bytes.Equal(uPasswordBytes, hashPassword), nil
}

// fillUserFields - method used to make finish view of updated user
func fillUserFields(uDest *user.User, uSrc *user.User) {
	if uSrc.GetName() != "" {
		uDest.Name = uSrc.GetName()
	}
	if uSrc.GetEmail() != "" {
		uDest.Email = uSrc.GetEmail()
	}
	uDest.Phone = uSrc.GetPhone()
}

func (uc *Layer) IsPassswordEquals(ctx context.Context, data *user.PasswordCheck) (*wrapperspb.BoolValue, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		uc.logger.Error(err.Error())
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return &wrapperspb.BoolValue{Value: false}, grpcerr.NewError(codes.NotFound, err.Error())
		}
		return &wrapperspb.BoolValue{Value: false}, grpcerr.NewError(codes.Internal, err.Error())
	}
	uPasswordBytes := []byte(u.GetPassword())

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, data.GetPassword())
	return &wrapperspb.BoolValue{Value: bytes.Equal(uPasswordBytes, hashPassword)}, nil
}

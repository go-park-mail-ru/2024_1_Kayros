package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/minios3"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/microservices/user/internal/repo"
	userv1 "2024_1_kayros/microservices/user/proto"

	"github.com/satori/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecase interface {
	userv1.UnsafeUserManagerServer
	GetData(ctx context.Context, email *userv1.Email) (*userv1.User, error)
	UpdateData(ctx context.Context, data *userv1.UpdateUserData) (*userv1.User, error)
	Create(ctx context.Context, data *userv1.User) (*userv1.User, error)
	UpdateAddress(ctx context.Context, data *userv1.AddressData) (*emptypb.Empty, error)
	SetNewPassword(ctx context.Context, data *userv1.PasswordsChange) (*emptypb.Empty, error)
	UpdateAddressByUnauthId(ctx context.Context, unauth *userv1.AddressDataUnauth) (*emptypb.Empty, error)
	GetAddressByUnauthId(ctx context.Context, id *userv1.UnauthId) (*userv1.Address, error)
}

type Layer struct {
	userv1.UnsafeUserManagerServer
	repoUser repo.Repo
	minio    minios3.Repo
}

func NewLayer(repoUserProps repo.Repo, repoMinio minios3.Repo) Usecase {
	return &Layer{
		repoUser: repoUserProps,
		minio:    repoMinio,
	}
}

// GetData - method calls repo method to receive user data.
func (uc Layer) GetData(ctx context.Context, email *userv1.Email) (*userv1.User, error) {
	return uc.repoUser.GetByEmail(ctx, email)
}

// UpdateData - method used to update user info. it accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc Layer) UpdateData(ctx context.Context, data *userv1.UpdateUserData) (*userv1.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		return nil, err
	}

	fillUserFields(u, data.GetUpdateInfo())
	if data.GetFileData() != nil && data.GetFileSize() != 0 {
		mimeType, err := functions.GetFileMimeType(data.GetFileData())
		if err != nil {
			return nil, err
		}
		if _, ok := cnst.ValidMimeTypes[mimeType]; !ok {
			return nil, myerrors.WrongFileExtension
		}

		fileExtension := functions.GetFileExtension(data.GetFileName())
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, bytes.NewBuffer(data.GetFileData()), filename, data.GetFileSize(), mimeType)
		if err != nil {
			return nil, err
		}
		u.ImgUrl = fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	}

	uDB, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if uDB != nil && data.GetEmail().GetEmail() != data.GetUpdateInfo().GetEmail().GetEmail() {
		return nil, myerrors.UserAlreadyExist
	}

	updateDataProps := &userv1.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateAddress - method updates only address.
func (uc Layer) UpdateAddress(ctx context.Context, data *userv1.AddressData) (*emptypb.Empty, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		return nil, err
	}

	u.Address = data.GetAddress()
	updateDataProps := &userv1.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc Layer) GetAddressByUnauthId(ctx context.Context, id *userv1.UnauthId) (*userv1.Address, error) {
	return uc.repoUser.GetAddressByUnauthId(ctx, id)
}

// SetNewPassword - method used to set new password.
func (uc Layer) SetNewPassword(ctx context.Context, data *userv1.PasswordsChange) (*emptypb.Empty,error) {
	// compare old password with database password
	isEqual, err := uc.checkPassword(ctx, data.GetEmail(), data.GetPassword())
	if err != nil {
		return nil, err
	}
	if !isEqual {
		return nil, myerrors.IncorrectCurrentPassword
	}
	if data.GetPassword() == data.GetNewPassword() {
		return nil, myerrors.NewPassword
	}

	u, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		return nil, err
	}

	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetNewPassword().GetPassword())
	u.Password = &userv1.Password{Password: string(hashPassword)}

	updateDataProps := &userv1.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (uc Layer) UpdateAddressByUnauthId(ctx context.Context, unauth *userv1.AddressDataUnauth) (*emptypb.Empty, error) {
	_, err := uc.repoUser.GetAddressByUnauthId(ctx, unauth.GetId())
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
			return nil, uc.repoUser.CreateAddressByUnauthId(ctx, unauth)
		}
		return nil, err
	}
	return nil, uc.repoUser.UpdateAddressByUnauthId(ctx, unauth)
}

// Create - method used to create new user in database.
func (uc Layer) Create(ctx context.Context, data *userv1.User) (*userv1.User, error) {
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetPassword().GetPassword())

	uCopy := entity.Copy(data)
	uCopy.Password = &userv1.Password{Password: string(hashPassword)}

	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}

	u, err := uc.repoUser.GetByEmail(ctx, uCopy.GetEmail())
	if err != nil {
		return nil, err
	}
	return u, nil
}

// checkPassword - method used to check password with password saved in database
func (uc Layer) checkPassword(ctx context.Context, email *userv1.Email, password *userv1.Password) (bool, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	uPasswordBytes := []byte(u.GetPassword().GetPassword())

	hashPassword := functions.HashData(uPasswordBytes[0:8], password.GetPassword())
	return bytes.Equal(uPasswordBytes, hashPassword), nil
}

// fillUserFields - method used to get finished view of updated user data
func fillUserFields(uDest *userv1.User, uSrc *userv1.User) {
	if uSrc.Name != "" {
		uDest.Name = uSrc.Name
	}
	if uSrc.Email.GetEmail() != "" {
		uDest.Email = uSrc.Email
	}
	uDest.Phone = uSrc.Phone
}

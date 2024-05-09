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
	"2024_1_kayros/gen/go/user"
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
}

func NewLayer(repoUserProps repo.Repo, repoMinio minios3.Repo) Usecase {
	return &Layer{
		repoUser: repoUserProps,
		minio:    repoMinio,
	}
}

func deleteCredentials (u *user.User) *user.User {
	if u != nil {
		u.Password = ""
		u.CardNumber = ""
	}
	return u
} 

// GetData - method calls repo method to receive user data.
func (uc Layer) GetData(ctx context.Context, email *user.Email) (*user.User, error) {
	returnUser, err := uc.repoUser.GetByEmail(ctx, email)
	return deleteCredentials(returnUser), err
}

// UpdateData - method used to update user info. it accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc Layer) UpdateData(ctx context.Context, data *user.UpdateUserData) (*user.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
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

	newEmail := &user.Email{Email: data.GetUpdateInfo().GetEmail()}
	uDB, err := uc.repoUser.GetByEmail(ctx, newEmail)
	if err != nil {
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if uDB != nil && data.GetEmail() != newEmail.GetEmail() {
		return nil, myerrors.UserAlreadyExist
	}

	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}
	return deleteCredentials(u), nil
}

// UpdateAddress - method updates only address.
func (uc Layer) UpdateAddress(ctx context.Context, data *user.AddressData) (*emptypb.Empty, error) {
	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		return nil, err
	}

	u.Address = data.GetAddress()
	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (uc Layer) GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error) {
	return uc.repoUser.GetAddressByUnauthId(ctx, id)
}

// SetNewPassword - method used to set new password.
func (uc Layer) SetNewPassword(ctx context.Context, data *user.PasswordsChange) (*emptypb.Empty,error) {
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

	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: data.GetEmail()})
	if err != nil {
		return nil, err
	}
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetNewPassword())
	u.Password = string(hashPassword)

	updateDataProps := &user.UpdateUserData{
		Email:      data.GetEmail(),
		UpdateInfo: u,
	}
	err = uc.repoUser.Update(ctx, updateDataProps)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (uc Layer) UpdateAddressByUnauthId(ctx context.Context, unauth *user.AddressDataUnauth) (*emptypb.Empty, error) {
	_, err := uc.repoUser.GetAddressByUnauthId(ctx, &user.UnauthId{UnauthId: unauth.GetUnauthId()})
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
			return nil, uc.repoUser.CreateAddressByUnauthId(ctx, unauth)
		}
		return nil, err
	}
	return nil, uc.repoUser.UpdateAddressByUnauthId(ctx, unauth)
}

// Create - method used to create new user in database.
func (uc Layer) Create(ctx context.Context, data *user.User) (*user.User, error) {
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, data.GetPassword())

	uCopy := entity.Copy(data)
	uCopy.Password = string(hashPassword)

	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}

	u, err := uc.repoUser.GetByEmail(ctx, &user.Email{Email: uCopy.GetEmail()})
	if err != nil {
		return nil, err
	}
	return deleteCredentials(u), nil
}

// checkPassword - method used to check password with password saved in database
func (uc Layer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
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
		return &wrapperspb.BoolValue{Value: false}, err
	}
	uPasswordBytes := []byte(u.GetPassword())

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, data.GetPassword())
	return &wrapperspb.BoolValue{Value: bytes.Equal(uPasswordBytes, hashPassword)}, nil
}
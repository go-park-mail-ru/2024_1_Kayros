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
	"google.golang.org/protobuf/types/known/wrapperspb"

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
	IsPassswordEquals(ctx context.Context, pwds *userv1.PasswordCheck) (*wrapperspb.BoolValue, error)
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

func formatUser (u *userv1.User) *userv1.User {
	if u != nil {
		u.Password = &userv1.Password{}
		u.CardNumber = ""
	}
	return u
} 

// GetData - method calls repo method to receive user data.
func (uc Layer) GetData(ctx context.Context, email *userv1.Email) (*userv1.User, error) {
	returnUser, err := uc.repoUser.GetByEmail(ctx, email)
	fmt.Printf("%v", returnUser)
	fmt.Printf("%v", err)
	return formatUser(returnUser), err
}

// UpdateData - method used to update user info. it accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc Layer) UpdateData(ctx context.Context, data *userv1.UpdateUserData) (*userv1.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		return nil, err
	}
	fmt.Printf("ДО:\n%v", u)
	fillUserFields(u, data.GetUpdateInfo())
	fmt.Printf("ОБНОВА:\n%v", data.GetUpdateInfo())
	fmt.Printf("ПОСЛЕ:\n%v", u)
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

	uDB, err := uc.repoUser.GetByEmail(ctx, data.GetUpdateInfo().GetEmail())
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
	return formatUser(u), nil
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
	fmt.Printf("%v", uCopy)

	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", err)

	u, err := uc.repoUser.GetByEmail(ctx, uCopy.GetEmail())
	fmt.Printf("%v", err)
	fmt.Printf("%v", u)
	if err != nil {
		return nil, err
	}
	return formatUser(u), nil
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
	if uSrc.GetName() != "" {
		uDest.Name = uSrc.GetName() 
	}
	if uSrc.GetEmail().GetEmail() != "" {
		uDest.GetEmail().Email = uSrc.GetEmail().GetEmail()
	}
	if uSrc.GetAddress().GetAddress() != "" {
		uDest.GetAddress().Address = uSrc.GetAddress().GetAddress()
	}
	uDest.Phone = uSrc.GetPhone()
}

func (uc *Layer) IsPassswordEquals(ctx context.Context, data *userv1.PasswordCheck) (*wrapperspb.BoolValue, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.GetEmail())
	if err != nil {
		return &wrapperspb.BoolValue{Value: false}, err
	}
	uPasswordBytes := []byte(u.GetPassword().GetPassword())

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, data.GetPwd().GetPassword())
	return &wrapperspb.BoolValue{Value: bytes.Equal(uPasswordBytes, hashPassword)}, nil
}
package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/minios3"
	"2024_1_kayros/internal/repository/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"github.com/satori/uuid"
)

type Usecase interface {
	GetData(ctx context.Context, email string) (*entity.User, error)
	UpdateAddressByUnauthId(ctx context.Context, unauthId string, address string) error
	GetAddressByUnauthId(ctx context.Context, unauthId string) (string, error)
	UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error)
	UpdateAddress(ctx context.Context, email string, address string) error
	SetNewPassword(ctx context.Context, email string, pwds *props.SetNewUserPasswordProps) error
}

type UsecaseLayer struct {
	repoUser user.Repo
	minio    minios3.Repo
}

func NewUsecaseLayer(repoUserProps user.Repo, repoMinio minios3.Repo) Usecase {
	return &UsecaseLayer{
		repoUser: repoUserProps,
		minio:    repoMinio,
	}
}

// GetData - method calls repo method to receive user data.
func (uc *UsecaseLayer) GetData(ctx context.Context, email string) (*entity.User, error) {
	return uc.repoUser.GetByEmail(ctx, email)
}

// UpdateData - method used to update user info. it accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc *UsecaseLayer) UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	fillUserFields(u, data.UserPropsUpdate)
	if data.File != nil && data.Handler != nil {
		buffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(buffer, data.File); err != nil {
			return nil, err
		}
		mimeType, err := functions.GetFileMimeType(buffer.Bytes())
		if err != nil {
			return nil, err
		}
		fileExtension := functions.GetFileExtension(data.Handler.Filename)
		if _, ok := cnst.ValidMimeTypes[mimeType]; !ok {
			return nil, myerrors.WrongFileExtension
		}

		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, buffer, filename, data.Handler.Size, mimeType)
		if err != nil {
			return nil, err
		}
		u.ImgUrl = fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	}

	uDB, err := uc.repoUser.GetByEmail(ctx, data.UserPropsUpdate.Email)
	if err != nil {
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if uDB != nil && data.Email != data.UserPropsUpdate.Email {
		return nil, myerrors.UserAlreadyExist
	}

	err = uc.repoUser.Update(ctx, u, data.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateAddress - method updates only address.
func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, email string, address string) error {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	u.Address = address
	err = uc.repoUser.Update(ctx, u, email)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsecaseLayer) GetAddressByUnauthId(ctx context.Context, unauthId string) (string, error) {
	return uc.repoUser.GetAddressByUnauthId(ctx, unauthId)
}

func (uc *UsecaseLayer) UpdateAddressByUnauthId(ctx context.Context, unauthId string, addressUpdate string) error {
	_, err := uc.repoUser.GetAddressByUnauthId(ctx, unauthId)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
			return uc.repoUser.CreateAddressByUnauthId(ctx, unauthId, addressUpdate)
		}
		return err
	}
	return uc.repoUser.UpdateAddressByUnauthId(ctx, unauthId, addressUpdate)
}

// SetNewPassword - method used to set new password.
func (uc *UsecaseLayer) SetNewPassword(ctx context.Context, email string, pwds *props.SetNewUserPasswordProps) error {
	// compare old password with database password
	isEqual, err := uc.checkPassword(ctx, email, pwds.Password)
	if err != nil {
		return err
	}
	if !isEqual {
		return myerrors.IncorrectCurrentPassword
	}
	if pwds.Password == pwds.PasswordNew {
		return myerrors.NewPassword
	}

	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return err
	}
	hashPassword := functions.HashData(salt, pwds.PasswordNew)
	u.Password = string(hashPassword)

	err = uc.repoUser.Update(ctx, u, email)
	if err != nil {
		return err
	}
	return nil
}

// Create - method used to create new user in database.
func (uc *UsecaseLayer) Create(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, uProps.Password)

	uCopy := entity.Copy(uProps)
	uCopy.Password = string(hashPassword)

	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}

	u, err := uc.repoUser.GetByEmail(ctx, uCopy.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// checkPassword - method used to check password with password saved in database
func (uc *UsecaseLayer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	uPasswordBytes := []byte(u.Password)

	hashPassword := functions.HashData(uPasswordBytes[0:8], password)
	return bytes.Equal(uPasswordBytes, hashPassword), nil
}

// fillUserFields - method used to get finished view of updated user data
func fillUserFields(uDest *entity.User, uSrc *entity.User) {
	if uSrc.Name != "" {
		uDest.Name = uSrc.Name
	}
	if uSrc.Email != "" {
		uDest.Email = uSrc.Email
	}
	uDest.Phone = uSrc.Phone
}

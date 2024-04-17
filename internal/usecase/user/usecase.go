package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"

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
	UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error)
	UpdateAddress(ctx context.Context, email string, address string) (*entity.User, error)
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
// to update password use method SetNewUserPassword.
func (uc *UsecaseLayer) UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	fillUserFields(u, data.UserPropsUpdate)
	if data.File != nil && data.Handler != nil {
		mimeType, err := functions.GetFileMimeType(data.File)
		if err != nil {
			return nil, err
		}
		fileExtension := functions.GetFileExtension(mimeType)
		if _, ok := cnst.ValidMimeTypes[fileExtension]; !ok {
			return nil, myerrors.WrongFileExtension
		}
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, data.File, filename, data.Handler.Size)
		if err != nil {
			return nil, err
		}
		u.ImgUrl = fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	}

	err = uc.repoUser.Update(ctx, u, data.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// UpdateAddress - method updates only address.
func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, email string, address string) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	u.Address = address
	err = uc.repoUser.Update(ctx, u, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// SetNewPassword - method used to set new password.
func (uc *UsecaseLayer) SetNewPassword(ctx context.Context, email string, pwds *props.SetNewUserPasswordProps) error {
	// compare old password with database password
	isEqual, err := uc.checkPassword(ctx, email, pwds.Password)
	if err != nil {
		return err
	}
	if !isEqual {
		return myerrors.Password
	}
	if pwds.Password == pwds.PasswordNew {
		return myerrors.NewPassword
	}

	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	salt := make([]byte, 8)
	_, err = rand.Read(salt)
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
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
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

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, password)
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

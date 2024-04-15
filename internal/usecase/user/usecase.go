package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"mime/multipart"

	"2024_1_kayros/internal/repository/minios3"
	"github.com/satori/uuid"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

type Usecase interface {
	GetById(ctx context.Context, userId alias.UserId) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId) error
	DeleteByEmail(ctx context.Context, email string) error

	IsExistById(ctx context.Context, userId alias.UserId) (bool, error)
	IsExistByEmail(ctx context.Context, email string) (bool, error)

	Create(ctx context.Context, uProps *entity.User) (*entity.User, error)
	Update(ctx context.Context, email string, uPropsUpdate *entity.User, file multipart.File, handler *multipart.FileHeader) (*entity.User, error)

	CheckPassword(ctx context.Context, email string, password string) (bool, error)
}

type UsecaseLayer struct {
	repoUser user.Repo
	minio    minios3.Repo
	logger   *zap.Logger
}

func NewUsecaseLayer(repoUserProps user.Repo, repoMinio minios3.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoUser: repoUserProps,
		minio:    repoMinio,
		logger:   loggerProps,
	}
}

func (uc *UsecaseLayer) GetById(ctx context.Context, userId alias.UserId) (*entity.User, error) {
	return uc.repoUser.GetById(ctx, userId)
}

func (uc *UsecaseLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.repoUser.GetByEmail(ctx, email)
}

func (uc *UsecaseLayer) DeleteById(ctx context.Context, userId alias.UserId) error {
	return uc.repoUser.DeleteById(ctx, userId)
}

func (uc *UsecaseLayer) DeleteByEmail(ctx context.Context, email string) error {
	return uc.repoUser.DeleteByEmail(ctx, email)
}

func (uc *UsecaseLayer) IsExistById(ctx context.Context, userId alias.UserId) (bool, error) {
	u, err := uc.repoUser.GetById(ctx, userId)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}
	return true, nil
}

func (uc *UsecaseLayer) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}
	return true, nil
}

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

func (uc *UsecaseLayer) Update(ctx context.Context, email string, uPropsUpdate *entity.User, file multipart.File, handler *multipart.FileHeader) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = fillUserFields(u, uPropsUpdate)
	if err != nil {
		return nil, err
	}

	if file != nil && handler != nil {
		fileExtension := functions.GetFileExtension(handler.Filename)
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, file, filename, handler.Size)
		if err != nil {
			return nil, err
		}
		u.ImgUrl = fmt.Sprintf("/minios3-api/%s/%s", cnst.BucketUser, filename)
	}

	err = uc.repoUser.Update(ctx, u, email)
	if err != nil {
		return nil, err
	}

	uData, err := uc.repoUser.GetByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	return uData, nil
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (uc *UsecaseLayer) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
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

func fillUserFields(uDest *entity.User, uSrc *entity.User) error {
	if uSrc.Name != "" {
		uDest.Name = uSrc.Name
	}

	if uSrc.Phone != "" {
		uDest.Phone = uSrc.Phone
	}

	if uSrc.Email != "" {
		uDest.Email = uSrc.Email
	}

	if uSrc.Address != "" {
		uDest.Address = uSrc.Address
	}

	if uSrc.Password != "" {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			return err
		}
		uDest.Password = string(functions.HashData(salt, uSrc.Password))
	}

	if uSrc.CardNumber != "" {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			return err
		}
		uDest.CardNumber = string(functions.HashData(salt, uSrc.CardNumber))
	}

	return nil
}

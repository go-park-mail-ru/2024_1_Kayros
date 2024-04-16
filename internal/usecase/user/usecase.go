package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"2024_1_kayros/internal/repository/minios3"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/satori/uuid"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/services/logger"
)

type Usecase interface {
	GetUserData(ctx context.Context, email string, requestId string, myLogger *logger.MyLogger) (*entity.User, error)
	UpdateUserData(ctx context.Context, email string, file multipart.File, handler *multipart.FileHeader, uPropsUpdate *entity.User, requestId string, myLogger *logger.MyLogger) (*entity.User, error)
	UpdateUserAddress(ctx context.Context, email string, address string, requestId string, myLogger *logger.MyLogger) (*entity.User, error)
	SetNewUserPassword(ctx context.Context, userId alias.UserId, requestId string, myLogger *logger.MyLogger) (*entity.User, error)
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

func (uc *UsecaseLayer) GetUserData(ctx context.Context, email string, requestId string, myLogger *logger.MyLogger) (*entity.User, error) {
	return uc.repoUser.GetByEmail(ctx, email, requestId, myLogger)
}

func (uc *UsecaseLayer) UpdateUserData(ctx context.Context, email string, file multipart.File, handler *multipart.FileHeader, uPropsUpdate *entity.User, requestId string, myLogger *logger.MyLogger) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId, myLogger)
	if err != nil {
		return nil, err
	}

	err = fillUserFields(u, uPropsUpdate)
	if err != nil {
		return nil, err
	}

	// нужно будет добавить функцию, возвращающую расширение файла по содержимому файла, а не по факт расширению
	if file != nil && handler != nil {
		fileExtension := functions.GetFileExtension(handler.Filename)
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.minio.UploadImageByEmail(ctx, file, filename, handler.Size)
		if err != nil {
			return nil, err
		}
		u.ImgUrl = fmt.Sprintf("/minios3-api/%s/%s", cnst.BucketUser, filename)
	}

	err = uc.repoUser.Update(ctx, u, email, requestId, myLogger)
	if err != nil {
		return nil, err
	}

	uData, err := uc.repoUser.GetByEmail(ctx, u.Email, requestId, myLogger)
	if err != nil {
		return nil, err
	}

	return uData, nil
}

func (uc *UsecaseLayer) UpdateUserAddress(ctx context.Context, email string, address string, requestId string, myLogger *logger.MyLogger) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId, myLogger)
	if err != nil {
		return nil, err
	}

	u.Address = address
	err = uc.repoUser.Update(ctx, u, email, requestId, myLogger)
	if err != nil {
		return nil, err
	}

	uDB, err := uc.repoUser.GetByEmail(ctx, email, requestId, myLogger)
	if err != nil {
		return nil, err
	}
	return uDB, nil
}

func (uc *UsecaseLayer) SetNewUserPassword(ctx context.Context, userId alias.UserId, requestId string, myLogger *logger.MyLogger) (*entity.User, error) {
	// сравниваем старый пароль с тем, что в базе
	isEqual, err := d.ucUser.CheckPassword(r.Context(), email, password.Password)
	if err != nil {
		return nil, err
	}
	// если они не совпадают
	if !isEqual {
		return nil, err
	}
	// они совпадают, значит мы можем поменять пароль пользователю
	// проверяем, что старый и новый пароль должны быть разными
	if password.Password == password.NewPassword {
		err = errors.New(myerrors.EqualPasswordsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.EqualPasswordsError, http.StatusBadRequest)
		return
	}
	_, err = d.ucUser.SetNewPassword(r.Context(), email, password.NewPassword)
	if err != nil {
		err = errors.New(myerrors.BadCredentialsError)
		functions.LogErrorResponse(d.logger, requestId, cnst.NameHandlerUpdatePassword, err, http.StatusBadRequest, cnst.DeliveryLayer)
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
}

//func (uc *UsecaseLayer) IsExistById(ctx context.Context, userId alias.UserId, requestId string, myLogger *logger.MyLogger) (bool, error) {
//	u, err := uc.repoUser.GetById(ctx, userId, requestId, myLogger)
//	if err != nil {
//		return false, err
//	}
//	if u == nil {
//		return false, nil
//	}
//	return true, nil
//}
//
//func (uc *UsecaseLayer) IsExistByEmail(ctx context.Context, email string, requestId string, myLogger *logger.MyLogger) (bool, error) {
//	u, err := uc.repoUser.GetByEmail(ctx, email, requestId, myLogger)
//	if err != nil {
//		return false, err
//	}
//	if u == nil {
//		return false, nil
//	}
//	return true, nil
//}

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

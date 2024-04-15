package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

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
	Update(ctx context.Context, email string, file multipart.File, handler *multipart.FileHeader, uProps *entity.User) (*entity.User, error)

	CheckPassword(ctx context.Context, email string, password string) (bool, error)
	SetNewPassword(ctx context.Context, email string, password string) (bool, error)
}

type UsecaseLayer struct {
	repoUser user.Repo
	logger   *zap.Logger
}

func NewUsecaseLayer(repoUserProps user.Repo, loggerProps *zap.Logger) Usecase {
	return &UsecaseLayer{
		repoUser: repoUserProps,
		logger:   loggerProps,
	}
}

func (uc *UsecaseLayer) GetById(ctx context.Context, userId alias.UserId) (*entity.User, error) {
	methodName := cnst.NameMethodGetUserById
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	u, err := uc.repoUser.GetById(ctx, userId, requestId)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return u, err
}

func (uc *UsecaseLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	methodName := cnst.NameMethodGetUserByEmail
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	fmt.Println("we are in getuser")
	u, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err == nil {
		fmt.Println("we get user")
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
		return u, nil
	} else {
		fmt.Println("we have truble with user", err)
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return u, err
}

func (uc *UsecaseLayer) DeleteById(ctx context.Context, userId alias.UserId) error {
	methodName := cnst.NameMethodDeleteUserById
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	err := uc.repoUser.DeleteById(ctx, userId, requestId)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return err
}

func (uc *UsecaseLayer) DeleteByEmail(ctx context.Context, email string) error {
	methodName := cnst.NameMethodDeleteUserByEmail
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	err := uc.repoUser.DeleteByEmail(ctx, email, requestId)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return err
}

func (uc *UsecaseLayer) IsExistById(ctx context.Context, userId alias.UserId) (bool, error) {
	methodName := cnst.NameMethodIsExistById
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	isExist, err := uc.repoUser.IsExistById(ctx, userId, requestId)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return isExist, err
}

func (uc *UsecaseLayer) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	methodName := cnst.NameMethodIsExistByEmail
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	isExist, err := uc.repoUser.IsExistByEmail(ctx, email, requestId)
	if err == nil {
		functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	} else {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
	}
	return isExist, err
}

func (uc *UsecaseLayer) Create(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	methodName := cnst.NameMethodCreateUser
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	hashPassword := functions.HashData(salt, uProps.Password)

	currentTime := time.Now().UTC()
	timeStr := currentTime.Format("2006-01-02 15:04:05-07:00")
	err = uc.repoUser.Create(ctx, uProps, hashPassword, timeStr, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}

	u, err := uc.repoUser.GetByEmail(ctx, uProps.Email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return u, nil
}

func (uc *UsecaseLayer) Update(ctx context.Context, email string, file multipart.File, handler *multipart.FileHeader, uPropsUpdate *entity.User) (*entity.User, error) {
	methodName := cnst.NameMethodUpdateUser
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	uPassword, err := uc.repoUser.GetHashedUserPassword(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	uCardNumber, err := uc.repoUser.GetHashedCardNumber(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}

	hashPassword := []byte{}
	if uPropsUpdate.Password == "" {
		hashPassword = uPassword
	} else {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
		hashPassword = functions.HashData(salt, uPropsUpdate.Password)
	}

	hashCardNumber := []byte{}
	if uPropsUpdate.CardNumber == "" {
		hashCardNumber = uCardNumber
	} else {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
		hashCardNumber = functions.HashData(salt, uPropsUpdate.CardNumber)
	}

	currentTime := time.Now().UTC()
	timeStr := currentTime.Format("2006-01-02 15:04:05-07:00")
	if file != nil && handler != nil {
		fileExtension := functions.GetFileExtension(handler.Filename)
		filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
		err = uc.repoUser.UploadImageByEmail(ctx, file, filename, handler.Size, email, timeStr, requestId)
		if err != nil {
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
	}
	uOldData, err := uc.repoUser.GetByEmail(ctx, email, requestId)
	if err != nil || uOldData == nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	if uPropsUpdate.Name == "" {
		return nil, errors.New("Имя не может быть пустым")
	}
	if uPropsUpdate.Email == "" {
		return nil, errors.New("Почта не может быть пустой")
	}
	if uPropsUpdate.ImgUrl == "" {
		uPropsUpdate.ImgUrl = uOldData.ImgUrl
	}
	if uPropsUpdate.Address == "" {
		uPropsUpdate.Address = uOldData.Address
	}

	err = uc.repoUser.Update(ctx, email, uPropsUpdate, hashPassword, hashCardNumber, timeStr, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	u, err := uc.repoUser.GetByEmail(ctx, uPropsUpdate.Email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return u, nil
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (uc *UsecaseLayer) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
	methodName := cnst.NameMethodCheckPassword
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	uPassword, err := uc.repoUser.GetHashedUserPassword(ctx, email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return false, err
	}

	salt := make([]byte, 8)
	copy(salt, uPassword[0:8])
	hashPassword := functions.HashData(salt, password)

	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return bytes.Equal(uPassword, hashPassword), nil
}

// SetNewPassword устанавливает новый пароль пользователю
func (uc *UsecaseLayer) SetNewPassword(ctx context.Context, email string, password string) (bool, error) {
	methodName := cnst.NameMethodSetNewPassword
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)

	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return false, err
	}
	hashPassword := functions.HashData(salt, password)
	_, err = uc.repoUser.SetNewPassword(ctx, requestId, email, hashPassword)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return false, err
	}
	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return true, nil
}

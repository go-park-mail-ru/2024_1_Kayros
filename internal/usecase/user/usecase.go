package user

import (
	"context"
	"crypto/rand"
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
	Update(ctx context.Context, uProps *entity.User) (*entity.User, error)

	CheckPassword(ctx context.Context, email string, password string) (bool, error)
	UploadImageByEmail(ctx context.Context, file multipart.File, handler *multipart.FileHeader, email string) error
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
	hashPasswordBytes := functions.HashData(salt, uProps.Password)
	hashPassword := string(hashPasswordBytes)

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

func (uc *UsecaseLayer) Update(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	methodName := cnst.NameMethodUpdateUser
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	uPassword, err := uc.repoUser.GetHashedUserPassword(ctx, uProps.Email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	uCardNumber, err := uc.repoUser.GetHashedCardNumber(ctx, uProps.Email, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}

	var hashPassword string
	if uProps.Password == "" {
		hashPassword = uPassword
	} else {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
		hashPasswordByte := functions.HashData(salt, uProps.Password)
		hashPassword = string(hashPasswordByte)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
	}

	var hashСardNumber string
	if uProps.CardNumber == "" {
		hashСardNumber = uCardNumber
	} else {
		salt := make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
		hashCardNumberByte := functions.HashData(salt, uProps.CardNumber)
		hashСardNumber = string(hashCardNumberByte)
		if err != nil {
			functions.LogError(uc.logger, requestId, methodName, err, cnst.UsecaseLayer)
			functions.LogUsecaseFail(uc.logger, requestId, methodName)
			return nil, err
		}
	}

	currentTime := time.Now().UTC()
	timeStr := currentTime.Format("2006-01-02 15:04:05-07:00")
	err = uc.repoUser.Update(ctx, uProps, hashPassword, hashСardNumber, timeStr, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return nil, err
	}
	u, err := uc.repoUser.GetById(ctx, alias.UserId(uProps.Id), requestId)
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

	salt := uPassword[0:8]
	hashPasswordBytes := functions.HashData([]byte(salt), password)
	hashPassword := string(hashPasswordBytes)

	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return uPassword == hashPassword, nil
}

func (uc *UsecaseLayer) UploadImageByEmail(ctx context.Context, file multipart.File, handler *multipart.FileHeader, email string) error {
	methodName := cnst.NameMethodCheckPassword
	requestId := functions.GetRequestId(ctx, uc.logger, methodName)
	fileExtension := functions.GetFileExtension(handler.Filename)
	filename := fmt.Sprintf("%s.%s", uuid.NewV4().String(), fileExtension)
	currentTime := time.Now().UTC()
	timeStr := currentTime.Format("2006-01-02 15:04:05-07:00")
	err := uc.repoUser.UploadImageByEmail(ctx, file, filename, handler.Size, email, timeStr, requestId)
	if err != nil {
		functions.LogUsecaseFail(uc.logger, requestId, methodName)
		return err
	}
	functions.LogOk(uc.logger, requestId, methodName, cnst.UsecaseLayer)
	return nil
}

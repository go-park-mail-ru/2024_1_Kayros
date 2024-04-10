package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

// Передаем контекст запроса пользователя (! возможно лучше еще переопределить контекстом WithTimeout)
type Repo interface {
	GetById(ctx context.Context, userId alias.UserId, requestId string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string, requestId string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId, requestId string) error
	DeleteByEmail(ctx context.Context, email string, requestId string) error

	Create(ctx context.Context, u *entity.User, hashPassword []byte, timeStr string, requestId string) error
	Update(ctx context.Context, email string, u *entity.User, hashPassword []byte, hashCardNumber []byte, timeStr string, requestId string) error

	IsExistById(ctx context.Context, userId alias.UserId, requestId string) (bool, error)
	IsExistByEmail(ctx context.Context, email string, requestId string) (bool, error)

	UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64, email string, timeStr string, requestId string) error
	GetHashedUserPassword(ctx context.Context, email string, requestId string) ([]byte, error)
	GetHashedCardNumber(ctx context.Context, email string, requestId string) ([]byte, error)
}

type RepoLayer struct {
	database *sql.DB
	minio    *minio.Client
	logger   *zap.Logger
}

func NewRepoLayer(db *sql.DB, minioProps *minio.Client, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		database: db,
		minio:    minioProps,
		logger:   loggerProps,
	}
}

func (repo *RepoLayer) GetById(ctx context.Context, userId alias.UserId, requestId string) (*entity.User, error) {
	methodName := cnst.NameMethodGetUserById
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(phone, ''), email, COALESCE(address, ''), img_url FROM "user" WHERE id = $1`, uint64(userId))
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Address, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Вернулось пустое множество данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был получен из базы данных", user.Id, user.Email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return &user, nil
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string, requestId string) (*entity.User, error) {
	methodName := cnst.NameMethodGetUserByEmail
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(phone, ''), email, COALESCE(address, ''), img_url FROM "user" WHERE email = $1`, email)
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Address, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Вернулось пустое множество данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был получен из базы данных", user.Id, user.Email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return &user, nil
}

func (repo *RepoLayer) DeleteById(ctx context.Context, userId alias.UserId, requestId string) error {
	methodName := cnst.NameMethodDeleteUserById
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "user" WHERE id = $1 RETURNING id, email`, uint64(userId))
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Пользователя с таким Id нет. Удалить не получилось")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был удален из базы данных", uId, uEmail)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string, requestId string) error {
	methodName := cnst.NameMethodDeleteUserByEmail
	row := repo.database.QueryRowContext(ctx, `DELETE FROM "user" WHERE email = $1 RETURNING id, email`, email)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Пользователя с таким Email нет. Удалить не получилось")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был удален из базы данных", uId, uEmail)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User, hashPassword []byte, timeStr string, requestId string) error {
	methodName := cnst.NameMethodCreateUser
	row := repo.database.QueryRowContext(ctx,
		`INSERT INTO "user" (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, email`,
		u.Name, u.Email, hashPassword, timeStr, timeStr)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Ошибка получения данных после их добавления в базу данных")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был добавлен в базу данных", uId, uEmail)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) Update(ctx context.Context, email string, uDataChange *entity.User, hashPassword []byte, hashCardNumber []byte, timeStr string, requestId string) error {
	methodName := cnst.NameMethodUpdateUser
	row := repo.database.QueryRowContext(ctx,
		`UPDATE "user" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5, card_number = $6, updated_at = $7 WHERE email = $8 RETURNING id, email`,
		uDataChange.Name, uDataChange.Phone, uDataChange.Email, uDataChange.ImgUrl, hashPassword, hashCardNumber, timeStr, email)
	var uId uint64
	var uEmail string
	err := row.Scan(&uId, &uEmail)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Ошибка получения данных после их обновления в базе данных")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s был обновлен", uId, uEmail)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) IsExistById(ctx context.Context, userId alias.UserId, requestId string) (bool, error) {
	methodName := cnst.NameMethodIsExistById
	u, err := repo.GetById(ctx, userId, requestId)
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, err
	}
	if u == nil {
		err = errors.New("Пользователя нет в базе данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, nil
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s существует", u.Id, u.Email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return true, nil
}

func (repo *RepoLayer) IsExistByEmail(ctx context.Context, email string, requestId string) (bool, error) {
	methodName := cnst.NameMethodIsExistByEmail
	u, err := repo.GetByEmail(ctx, email, requestId)
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, err
	}
	if u == nil {
		err = errors.New("Пользователя нет в базе данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return false, nil
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s существует", u.Id, u.Email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return true, nil
}

func (repo *RepoLayer) GetHashedUserPassword(ctx context.Context, email string, requestId string) ([]byte, error) {
	methodName := cnst.NameMethodGetHashedUserPassword
	row := repo.database.QueryRowContext(ctx, `SELECT id, password FROM "user" WHERE email = $1`, email)
	var uId uint64
	var hashedPassword []byte
	err := row.Scan(&uId, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Пользователя с таким Email нет. Получить пароль не вышло")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return []byte{}, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return []byte{}, err
	}
	msg := fmt.Sprintf("Получили пароль пользователя с идентификатором %d и почтой %s", uId, email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return hashedPassword, nil
}

func (repo *RepoLayer) GetHashedCardNumber(ctx context.Context, email string, requestId string) ([]byte, error) {
	methodName := cnst.NameMethodGetHashedUserPassword
	row := repo.database.QueryRowContext(ctx, `SELECT id, card_number FROM "user" WHERE email = $1`, email)
	var uId uint64
	var hashedCardNumber []byte
	err := row.Scan(&uId, &hashedCardNumber)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Пользователя с таким Email нет. Получить пароль не вышло")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return []byte{}, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return []byte{}, err
	}
	msg := fmt.Sprintf("Получили пароль пользователя с идентификатором %d и почтой %s", uId, email)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return hashedCardNumber, nil
}

func (repo *RepoLayer) UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64, email string, timeStr string, requestId string) error {
	methodName := cnst.NameMethodUploadImageByEmail
	_, err := repo.minio.PutObject(ctx, cnst.BucketUser, filename, file, filesize, minio.PutObjectOptions{ContentType: "application/form-data"})
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}

	imgPath := fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	row := repo.database.QueryRowContext(ctx, `UPDATE "user" SET img_url = $1, updated_at = $2 WHERE email = $3 RETURNING id, email, img_url`, imgPath, timeStr, email)
	var uId uint64
	var uEmail string
	var uImg string
	err = row.Scan(&uId, &uEmail, &uImg)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Ошибка получения данных после их обновления в базе данных")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	msg := fmt.Sprintf("Пользователь с идентификатором %d и почтой %s имеет фото по адресу %s", uId, uEmail, uImg)
	functions.LogInfo(repo.logger, requestId, methodName, msg, cnst.RepoLayer)
	return nil
}

package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// Передаем контекст запроса пользователя (! возможно лучше еще переопределить контекстом WithTimeout)
type Repo interface {
	GetById(ctx context.Context, userId alias.UserId, requestId string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string, requestId string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId, requestId string) error
	DeleteByEmail(ctx context.Context, email string, requestId string) error

	Create(ctx context.Context, u *entity.User, hashPassword string, requestId string) error
	Update(ctx context.Context, u *entity.User, hashPassword string, requestId string) error

	IsExistById(ctx context.Context, userId alias.UserId, requestId string) (bool, error)
	IsExistByEmail(ctx context.Context, email string, requestId string) (bool, error)

	UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64, email string, requestId string) error
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
		`SELECT id, name, password, phone, email, address, img_url FROM "User" WHERE id = $1`, uint64(userId))
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Вернулось пустое множество данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, err
	}

	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return &user, nil
}

func (repo *RepoLayer) GetByEmail(ctx context.Context, email string, requestId string) (*entity.User, error) {
	methodName := cnst.NameMethodGetUserByEmail
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, password, phone, email, address, img_url FROM "User" WHERE email = $1`, email)
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.Password, &user.ImgUrl)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Вернулось пустое множество данных")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, nil
	}
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return nil, err
	}

	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return &user, nil
}

func (repo *RepoLayer) DeleteById(ctx context.Context, userId alias.UserId, requestId string) error {
	methodName := cnst.NameMethodDeleteUserById
	res, err := repo.database.ExecContext(ctx, `DELETE FROM "User" WHERE id = $1`, uint64(userId))
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if countRows == 0 {
		err = errors.New("Пользователь не был удален")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) DeleteByEmail(ctx context.Context, email string, requestId string) error {
	methodName := cnst.NameMethodDeleteUserByEmail
	res, err := repo.database.ExecContext(ctx, `DELETE FROM "User" WHERE email = $1`, email)
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if countRows == 0 {
		err = errors.New("Пользователь не был удален")
		functions.LogWarn(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) Create(ctx context.Context, u *entity.User, hashPassword string, requestId string) error {
	methodName := cnst.NameMethodCreateUser
	res, err := repo.database.ExecContext(ctx,
		`INSERT INTO "User" (name, phone, email, password, img_url) VALUES ($1, $2, $3, $4, $5)`,
		u.Name, u.Phone, u.Email, hashPassword, u.ImgUrl)
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if countRows == 0 {
		err = errors.New("Пользователь не был добавлен")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return nil
}

func (repo *RepoLayer) Update(ctx context.Context, u *entity.User, hashPassword string, requestId string) error {
	methodName := cnst.NameMethodUpdateUser
	res, err := repo.database.ExecContext(ctx,
		`UPDATE "User" SET name = $1, phone = $2, email = $3, img_url = $4, password = $5 WHERE id = $6`,
		u.Name, u.Phone, u.Email, u.ImgUrl, hashPassword, u.Id)

	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if countRows == 0 {
		err = errors.New("Данные о пользователе не были обновлены")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
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
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
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
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return true, nil
}

func (repo *RepoLayer) UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64, email string, requestId string) error {
	methodName := cnst.NameMethodUploadImageByEmail
	_, err := repo.minio.PutObject(ctx, cnst.BucketUser, filename, file, filesize, minio.PutObjectOptions{ContentType: "application/form-data"})
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}

	imgPath := fmt.Sprintf("/minio-api/%s/%s", cnst.BucketUser, filename)
	res, err := repo.database.ExecContext(ctx, `UPDATE "User" SET img_url = $1 WHERE email = $2`, imgPath, email)
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	if countRows == 0 {
		err = errors.New("Фотография пользователя не была добавлена")
		functions.LogError(repo.logger, requestId, methodName, err, cnst.RepoLayer)
		return err
	}
	functions.LogOk(repo.logger, requestId, methodName, cnst.RepoLayer)
	return nil
}

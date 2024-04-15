package minios3

import (
	"context"
	"mime/multipart"

	cnst "2024_1_kayros/internal/utils/constants"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Repo interface {
	UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64) error
}

type RepoLayer struct {
	minio  *minio.Client
	logger *zap.Logger
}

func NewRepoLayer(minioClient *minio.Client, loggerProps *zap.Logger) Repo {
	return &RepoLayer{
		minio:  minioClient,
		logger: loggerProps,
	}
}

func (repo *RepoLayer) UploadImageByEmail(ctx context.Context, file multipart.File, filename string, filesize int64) error {
	_, err := repo.minio.PutObject(ctx, cnst.BucketUser, filename, file, filesize, minio.PutObjectOptions{ContentType: "application/form-data"})
	return err
}

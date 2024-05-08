package minio

import (
	"context"
	"fmt"
	"time"

	"2024_1_kayros/config"
	cnst "2024_1_kayros/internal/utils/constants"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func Init(cfg *config.Project, logger *zap.Logger) *minio.Client {
	client, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.SslMode,
	})
	if err != nil {
		logger.Fatal("Failed to connect to Minio", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	buckets := []string{cnst.BucketUser, cnst.BucketRest, cnst.BucketFood}
	for _, bucket := range buckets {
		makeBucket(client, bucket, ctx, logger)
	}

	logger.Info("Minio connected successfully")
	return client
}

func makeBucket(client *minio.Client, bucket string, ctx context.Context, logger *zap.Logger) {
	err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		isExist, err := client.BucketExists(ctx, bucket)
		if err == nil && isExist {
			msg := fmt.Sprintf("A bucket named %s already exists", bucket)
			logger.Info(msg)
			return
		} else {
			msg := fmt.Sprintf("Creating a bucket with a name %s failed", bucket)
			logger.Fatal(msg, zap.Error(err))
		}
	}
	msg := fmt.Sprintf("A bucker named %s created successfully", bucket)
	logger.Info(msg)
}

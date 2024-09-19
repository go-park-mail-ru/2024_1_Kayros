package minio

import (
	"context"
	"fmt"
	"time"

	cfg "2024_1_kayros/config"
	cnst "2024_1_kayros/internal/utils/constants"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *minio.Client {
	minioCfg := cfg.Config.Minio
	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.AccessKey, minioCfg.SecretKey, ""),
		Secure: minioCfg.SslMode,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to connect to Minio: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	buckets := []string{cnst.BucketUser, cnst.BucketRest, cnst.BucketFood}
	for _, bucket := range buckets {
		makeBucket(client, bucket, ctx, logger)
	}
	logger.Info("minio connected successfully")
	return client
}

func makeBucket(client *minio.Client, bucket string, ctx context.Context, logger *zap.Logger) {
	err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		isExist, err := client.BucketExists(ctx, bucket)
		if err == nil && isExist {
			logger.Info(fmt.Sprintf("a bucket with a name '%s' already exists", bucket))
			return
		} else {
			logger.Fatal(fmt.Sprintf("creating a bucket with a name %s failed: %v", bucket, err))
		}
	}
	logger.Info(fmt.Sprintf("a bucket named %s created successfully", bucket))
}

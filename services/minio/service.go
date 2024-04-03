package minio

import (
	"context"
	"log"
	"sync"
	"time"

	"2024_1_kayros/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Init(cfg *config.Project) (*minio.Client, error) {
	ctxInit, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.SslMode,
	})
	if err != nil {
		log.Fatalln(err)
	}

	location := "photos"
	buckets := []string{"restaurants", "users", "foods"}
	var wgBucket *sync.WaitGroup
	wgBucket.Add(len(buckets))
	for _, bucket := range buckets {
		go makeBucket(client, bucket, ctxInit, location, wgBucket)
	}
	wgBucket.Wait()

	return client, nil
}

func makeBucket(client *minio.Client, bucket string, ctx context.Context, location string, wg *sync.WaitGroup) {
	err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucket)
	}
	wg.Done()
}

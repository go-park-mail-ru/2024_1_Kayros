package services

import (
	minios3 "2024_1_kayros/services/minio"
	"2024_1_kayros/services/postgres"
	"database/sql"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type Cluster struct {
	PsqlClient  *sql.DB
	MinioClient *minio.Client
}

func Init(logger *zap.Logger) *Cluster {
	return &Cluster{
		PsqlClient:  postgres.Init(logger),
		MinioClient: minios3.Init(logger),
	}
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"2024_1_kayros/config"
)

func Init(cfg *config.Project, logger *zap.Logger) *sql.DB {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password, dbConfig.SslMode)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to connect to PostgreSQL %s at address %s:%d\n",
			cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)
		logger.Fatal(errorMsg, zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wasConnected := false

	for i := 0; i < 3; i++ {
		err = db.PingContext(ctx)
		if err != nil {
			logger.Error("Test queries to PostgreSQL failed", zap.String("error", err.Error()))
		} else {
			wasConnected = true
			break
		}
		time.Sleep(3 * time.Second)
	}
	if !wasConnected {
		logger.Fatal("Unable to connect to PostgreSQL", zap.String("error", err.Error()))
	}

	// maximum number of open connections
	db.SetMaxOpenConns(int(dbConfig.MaxOpenConns))
	// maximum amount of time the connection can be reused
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)
	// maximum number of connections in the pool of idle connections
	db.SetMaxIdleConns(int(dbConfig.MaxIdleConns))
	// maximum time during which the connection can be idle
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime) * time.Second)

	logger.Info("PostgreSQL connected successfully")
	return db
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	cfg "2024_1_kayros/config"
)

func Init(logger *zap.Logger) *sql.DB {
	dbConfig := cfg.Config.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password, dbConfig.SslMode)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to connect to psql at address %s:%d\n: %v",
			dbConfig.Host, dbConfig.Port, err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	wasConnected := false

	for i := 0; i < 3; i++ {
		err = db.PingContext(ctx)
		if err == nil {
			wasConnected = true
			break
		}
		logger.Warn(fmt.Sprintf("test queries to psql failed: %v", err))
		time.Sleep(2 * time.Second)
	}
	if !wasConnected {
		logger.Fatal(fmt.Sprintf("unable to connect to psql: %v", err))
	}

	// maximum number of open connections
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// maximum amount of time the connection can be reused
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	// maximum number of connections in the pool of idle connections
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	// maximum time during which the connection can be idle
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	logger.Info("psql connected successfully")
	return db
}

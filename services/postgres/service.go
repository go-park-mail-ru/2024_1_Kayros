package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"2024_1_kayros/config"
)

// DatabaseInit инициализирует коннект с базой данных
func DatabaseInit(cfg *config.Project) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.SslMode, dbConfig.Host, dbConfig.Port)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(2 * time.Second)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	return db, nil
}

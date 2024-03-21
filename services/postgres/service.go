package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"2024_1_kayros/config"
)

// PostgresInit инициализирует коннект с базой данных
func PostgresInit(cfg *config.Project) (*sql.DB, error) {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password, dbConfig.SslMode)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// максимальное количество открытых соединений
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// максимальное количество времени, в течение которого соединение может быть повторно использовано
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime * time.Second)
	// максимальное количество соединений в пуле простаивающих соединений
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	// максимальное время, в течение которого соединение может быть бездействующим
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime * time.Second)

	return db, nil
}

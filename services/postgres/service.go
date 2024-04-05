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

// Init инициализирует коннект с базой данных
func Init(cfg *config.Project, logger *zap.Logger) *sql.DB {
	dbConfig := cfg.Postgres
	dataConnection := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password, dbConfig.SslMode)
	db, err := sql.Open("postgres", dataConnection)
	if err != nil {
		errorMsg := fmt.Sprintf("Не удалось подключиться к PostgreSQL %s по адресу %s:%d\n",
			cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)
		logger.Fatal(errorMsg, zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		logger.Fatal("Проверочные запрос к PostgreSQL завершились неудачно", zap.Error(err))
	}
	// максимальное количество открытых соединений
	db.SetMaxOpenConns(int(dbConfig.MaxOpenConns))
	// максимальное количество времени, в течение которого соединение может быть повторно использовано
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime * time.Second)
	// максимальное количество соединений в пуле простаивающих соединений
	db.SetMaxIdleConns(int(dbConfig.MaxIdleConns))
	// максимальное время, в течение которого соединение может быть бездействующим
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime * time.Second)

	logger.Info("PostgreSQL успешно подключен")
	return db
}

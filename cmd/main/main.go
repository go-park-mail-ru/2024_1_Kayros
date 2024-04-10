package main

import (
	"flag"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/app"
	"go.uber.org/zap"
)

func main() {
	var loggerLevel string
	flag.StringVar(&loggerLevel, "log_level", "prod", "Уровень логирования кода")
	flag.Parse()

	var logger *zap.Logger
	switch loggerLevel {
	case "dev":
		logger = zap.Must(zap.NewDevelopment())
	case "prod":
		logger = zap.Must(zap.NewProduction())
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Очистка данных после логгера прошла с ошибкой", zap.Error(err))
		}
	}(logger)

	cfg := config.NewConfig(logger)
	app.Run(cfg, logger)
	logger.Info("Сервер завершил работу")
}

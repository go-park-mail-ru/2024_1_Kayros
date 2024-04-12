package main

import (
	"2024_1_kayros/config"
	"2024_1_kayros/internal/app"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Info("Failed to sync logs into storage")
		}
	}(logger)

	cfg := config.NewConfig(logger)
	app.Run(cfg)
	logger.Info("The server has shut down")
}

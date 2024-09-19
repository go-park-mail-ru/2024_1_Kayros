package main

import (
	"2024_1_kayros/config"
	"2024_1_kayros/internal/app"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Warn(fmt.Sprintf("failed to sync logs into storage: %v", err))
		}
	}()
	config.Read(logger)
	app.Run(logger)
	logger.Info("the server has shut down")
}

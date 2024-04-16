package logger

import (
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type MyLogger struct {
	logger    *zap.Logger
	requestId string
}

func NewMyLogger() *MyLogger {
	return &MyLogger{
		logger:    zap.Must(zap.NewProduction()),
		requestId: uuid.NewV4().String(),
	}
}

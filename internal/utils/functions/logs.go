package functions

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"2024_1_kayros/internal/utils/constants"
)

func LogInfo(logger *zap.Logger, requestId string, methodName string, err error, layer string) {
	errorMsg := fmt.Sprintf("Запрос %s. Информация: %v", requestId, err.Error())
	logger.Info(errorMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", layer),
	)
}

func LogError(logger *zap.Logger, requestId string, methodName string, err error, layer string) {
	errorMsg := fmt.Sprintf("Запрос %s. Ошибка: %v", requestId, err.Error())
	logger.Error(errorMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", layer),
	)
}

func LogErrorResponse(logger *zap.Logger, requestId string, methodName string, err error, status int, layer string) {
	errorMsg := fmt.Sprintf("Запрос %s. Ошибка: %v", requestId, err.Error())
	logger.Error(errorMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.Int("response_status", status),
		zap.String("layer", layer),
	)
}

func LogWarn(logger *zap.Logger, requestId string, methodName string, err error, layer string) {
	errorMsg := fmt.Sprintf("Запрос %s. Предупреждение: %v", requestId, err.Error())
	logger.Warn(errorMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", layer),
	)
}

func LogOk(logger *zap.Logger, requestId string, methodName string, layer string) {
	infoMsg := fmt.Sprintf("Запрос %s. Успешно выполнился", requestId)
	logger.Info(infoMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", layer),
	)
}

func LogOkResponse(logger *zap.Logger, requestId string, methodName string, layer string) {
	infoMsg := fmt.Sprintf("Запрос %s. Успешно выполнился", requestId)
	logger.Info(infoMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.Int("response_status", http.StatusOK),
		zap.String("layer", layer),
	)
}

func LogUsecaseFail(logger *zap.Logger, requestId string, methodName string) {
	infoMsg := fmt.Sprintf("Запрос %s. Завершился с ошибкой.", requestId)
	logger.Error(infoMsg,
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", constants.UsecaseLayer),
	)
}

func LogNoRequestId(logger *zap.Logger, requestId string, methodName string) {
	logger.Error("Через контекст не был передан request_id",
		zap.String("method", methodName),
		zap.String("request_id", requestId),
		zap.String("layer", constants.UsecaseLayer),
	)
}

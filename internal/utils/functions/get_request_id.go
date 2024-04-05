package functions

import (
	"context"

	"go.uber.org/zap"
)

func GetRequestId(ctx context.Context, logger *zap.Logger, methodName string) string {
	requestId := ctx.Value("request_id").(string)
	if requestId == "" {
		LogNoRequestId(logger, requestId, methodName)
	}
	return requestId
}

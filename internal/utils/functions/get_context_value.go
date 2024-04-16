package functions

import (
	"errors"
	"net/http"

	"2024_1_kayros/services/logger"
)

func GetCtxLogger(r *http.Request) (*logger.MyLogger, error) {
	ctxLogger := r.Context().Value("logger")
	if ctxLogger == nil {
		err := errors.New("logger was not passed in the context")
		return nil, err
	}
	return ctxLogger.(*logger.MyLogger), nil
}

func GetCtxRequestId(r *http.Request) (string, error) {
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		err := errors.New("request_id was not passed in the context")
		return "", err
	}
	return ctxRequestId.(string), nil
}

func GetCtxEmail(r *http.Request) (string, error) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		err := errors.New("email was not passed in the context")
		return "", err
	}
	return ctxEmail.(string), nil
}

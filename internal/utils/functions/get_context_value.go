package functions

import (
	"net/http"

	"2024_1_kayros/internal/utils/myerrors"
)

func GetCtxRequestId(r *http.Request) (string, error) {
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		return "", myerrors.CtxRequestId
	}
	return ctxRequestId.(string), nil
}

func GetCtxEmail(r *http.Request) (string, error) {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		return "", myerrors.CtxEmail
	}
	return ctxEmail.(string), nil
}

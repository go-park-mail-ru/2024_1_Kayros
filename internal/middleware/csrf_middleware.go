package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// я хочу:
// хранить данные

// CsrfMiddleware
func CsrfMiddleware(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := ""
		emailCtx := r.Context().Value("email")
		if emailCtx != nil {
			email = emailCtx.(string)
		}
		requestId := ""
		requestIdCtx := r.Context().Value("request_id")
		if requestIdCtx != nil {
			requestId = requestIdCtx.(string)
		}

		handler.ServeHTTP(w, r)
	})
}

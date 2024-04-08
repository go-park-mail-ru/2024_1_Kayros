package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// CorsMiddleware решает политику SOP с помощью CORS-заголовков
func CorsMiddleware(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		logger.Info("Политики SOP разрешена")
		handler.ServeHTTP(w, r)
	})
}

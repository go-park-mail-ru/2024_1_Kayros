package middleware

import (
	"context"
	"net/http"
	"time"

	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

// SessionAuthentication добавляет в контекст почту пользователя, которого получилось аутентифицировать
func SessionAuthenticationMiddleware(handler http.Handler, ucUser user.Usecase, ucSession session.Usecase, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4().String()
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		ctxData, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		email, err := ucSession.GetValue(ctxData, alias.SessionKey(sessionCookie.Value))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		_, err = ucUser.GetByEmail(ctxData, string(email))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		ctx = context.WithValue(ctx, "request_id", requestId)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

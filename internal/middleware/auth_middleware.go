package middleware

import (
	"context"
	"errors"
	"net/http"

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
		sessionId := ""
		if sessionCookie != nil {
			sessionId = sessionCookie.Value
		}

		if errors.Is(err, http.ErrNoCookie) {
			functions.LogInfo(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
		} else if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		ctx := context.WithValue(r.Context(), "request_id", requestId)

		email, err := ucSession.GetValue(ctx, alias.SessionKey(sessionId))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		_, err = ucUser.GetByEmail(ctx, string(email))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			return
		}

		ctx = context.WithValue(ctx, "email", email)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

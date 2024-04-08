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

// SessionAuthenticationMiddleware добавляет в контекст email пользователя, которого получилось аутентифицировать, а также request_id
func SessionAuthenticationMiddleware(handler http.Handler, ucUser user.Usecase, ucSession session.Usecase, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4().String()
		sessionCookie, err := r.Cookie("session_id")
		sessionId := ""
		if sessionCookie != nil {
			sessionId = sessionCookie.Value
		}

		if errors.Is(err, http.ErrNoCookie) {
			functions.LogInfo(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err.Error(), cnst.MiddlewareLayer)
		} else if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			handler.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "request_id", requestId)

		email, err := ucSession.GetValue(ctx, alias.SessionKey(sessionId))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			handler.ServeHTTP(w, r)
			return
		}

		u, err := ucUser.GetByEmail(ctx, string(email))
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameSessionAuthenticationMiddleware, err, cnst.MiddlewareLayer)
			handler.ServeHTTP(w, r)
			return
		}

		if u == nil {
			email = ""
		}
		ctx = context.WithValue(ctx, "email", string(email))
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

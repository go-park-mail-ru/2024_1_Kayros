package middleware

import (
	"context"
	"errors"
	"net/http"

	"2024_1_kayros/internal/utils/myerrors"
	"go.uber.org/zap"

	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

// SessionAuthentication needed for authentication (check session cookie and if it exists, return associated user email)
func SessionAuthentication(handler http.Handler, ucUser user.Repo, ucSession session.Usecase, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := ""
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			logger.Warn(err.Error())
		}
		if sessionCookie != nil {
			sessionId = sessionCookie.Value
		}
		requestId, err := functions.GetCtxRequestId(r)
		if err != nil {
			logger.Error(err.Error())
		}

		ctx := r.Context()
		email, err := ucSession.GetValue(ctx, alias.SessionKey(sessionId))
		if err != nil {
			if !errors.Is(err, myerrors.RedisNoData) {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		}
		u, err := ucUser.GetByEmail(ctx, string(email))
		if err != nil {
			if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		}
		if u == nil {
			email = ""
		}

		ctx = context.WithValue(ctx, "email", string(email))
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

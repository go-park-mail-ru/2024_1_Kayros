package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/utils/myerrors"

	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

// SessionAuthentication needed for authentication (check session cookie and if it exists, return associated user email)
func SessionAuthentication(handler http.Handler, ucUser user.Usecase, ucSession session.Usecase, logger *zap.Logger, cfg *config.Redis) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := functions.GetCtxRequestId(r)
		sessionId, err := functions.GetCookieSessionValue(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}

		fmt.Printf("SessionId: %s", sessionId)
		fmt.Printf("cfg.DatabaseSession: %s", cfg.DatabaseSession)
		fmt.Printf("cfg.DatabaseCsrf: %s", cfg.DatabaseCsrf)
		email, err := ucSession.GetValue(r.Context(), alias.SessionKey(sessionId), int32(cfg.DatabaseSession))
		fmt.Printf("Email: %s", email)
		fmt.Printf("err: %s", err)
		if err != nil {
			if !errors.Is(err, myerrors.RedisNoData) {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
		}

		u, err := ucUser.GetData(r.Context(), string(email))
		if err != nil {
			if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
		}

		if u == nil {
			email = ""
		}
		ctx := context.WithValue(r.Context(), "email", string(email))
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

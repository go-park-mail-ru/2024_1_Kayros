package middleware

import (
	"context"
	"net/http"
	"time"

	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
)

// SessionAuthentication добавляет в контекст почту пользователя, которого получилось аутентифицировать
func SessionAuthenticationMiddleware(handler http.Handler, ucUser user.Usecase, ucSession session.Usecase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			return
		}

		ctxData, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		email, err := ucSession.GetValue(ctxData, alias.SessionKey(sessionCookie.Value))
		if err != nil {
			return
		}

		_, err = ucUser.GetByEmail(ctxData, string(email))
		if err != nil {
			return
		}

		var ctx context.Context
		ctx = context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

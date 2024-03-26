package middleware

import (
	"context"
	"net/http"

	repoSession "2024_1_kayros/internal/repository/session"
	repoUser "2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
)

// SessionAuthentication добавляет в контекст ключ авторизации пользователя, которого получилось аутентифицировать
func SessionAuthentication(m http.Handler, userRepo repoUser.UserRepositoryInterface, sessionRepo repoSession.SessionRepositoryInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			return
		}

		sessionId := sessionCookie.Value
		email, err := sessionRepo.GetValue(alias.SessionKey(sessionId))
		if err != nil {
			return
		}

		_, err = userRepo.GetByEmail(string(email))
		if err != nil {
			return
		}

		var ctx context.Context
		ctx = context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)

		m.ServeHTTP(w, r)
	})
}

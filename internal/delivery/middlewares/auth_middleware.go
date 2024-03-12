package middlewares

import (
	"context"
	"errors"
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

// SessionAuthentication добавляет в контекст ключ авторизации пользователя, которого получилось аутентифицировать
func SessionAuthentication(_ *mux.Router, db *entity.SystemDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if !errors.Is(errNoSessionCookie, http.ErrNoCookie) {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId == nil {
				// проверка на наличие UUID в таблице сессий
				userEmail, errGettingEmail := db.Sessions.GetValue(sessionId)
				if errGettingEmail == nil {
					_, errWrongCredentionals := db.Users.GetUser(userEmail)
					if errWrongCredentionals == nil {
						var ctx context.Context
						ctx = context.WithValue(r.Context(), "userKey", userEmail)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		var handler http.Handler
		handler.ServeHTTP(w, r)
	})
}

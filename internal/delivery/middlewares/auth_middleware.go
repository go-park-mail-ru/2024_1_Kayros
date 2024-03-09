package middlewares

import (
	"context"
	"errors"
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

// SessionAuthentication добавляет во временное хранилище Redis данные о пользователе, сделавшем запрос
// Если что-то неправ
func SessionAuthentication(handler http.Handler, db *entity.SystemDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if !errors.Is(errNoSessionCookie, http.ErrNoCookie) {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId == nil {
				// проверка на наличие UUID в таблице сессий
				userEmail, sessionNotExist := db.Sessions.GetValue(sessionId)
				if sessionNotExist == nil {
					db.Users.UsersMutex.RLock()
					user := Users[userEmail]
					UsersMu.RUnlock()

					var ctx context.Context
					ctx = context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				}
			}
		}
		handler.ServeHTTP(w, r)
	})
}

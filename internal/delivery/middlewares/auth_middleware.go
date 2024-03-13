package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

// SessionAuthentication добавляет в контекст ключ авторизации пользователя, которого получилось аутентифицировать
func SessionAuthentication(m http.Handler, db *entity.AuthDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if !errors.Is(errNoSessionCookie, http.ErrNoCookie) {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId == nil {
				// проверка на наличие сессии пользователя в таблице сессий
				key, errGettingKey := db.Sessions.GetValueByKey(sessionId)
				if errGettingKey == nil {
					_, errWrongCredentionals := db.Users.GetUser(key)
					if errWrongCredentionals == nil {
						var ctx context.Context
						log.Println("auth")
						log.Println(key)
						ctx = context.WithValue(r.Context(), "authKey", key)
						r = r.WithContext(ctx)
					}
				}
			}
		} else {
			log.Println("cookie is empty")
		}
		m.ServeHTTP(w, r)
	})
}

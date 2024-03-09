package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

// SessionAuthentication добавляет во временное хранилище Redis данные о пользователе, сделавшем запрос
func SessionAuthentication(handler http.Handler, db *entity.SystemDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if errors.Is(errNoSessionCookie, http.ErrNoCookie) {
			log.Println("Авторизационные Cookie пустые")
		} else {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId != nil {
				log.Println("Авторизационные Cookie имеют неверный формат")
			} else {
				// проверка на наличие UUID в таблице сессий
				userEmail, errGettingEmail := db.Sessions.GetValue(sessionId)
				if errGettingEmail != nil {
					log.Println(errGettingEmail)
				} else {
					user, errWrongCredentionals := db.Users.GetUser(userEmail)
					if errWrongCredentionals != nil {
						log.Println(errWrongCredentionals)
					} else {
						var ctx context.Context
						ctx = context.WithValue(r.Context(), "user", user)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		handler.ServeHTTP(w, r)
	})
}

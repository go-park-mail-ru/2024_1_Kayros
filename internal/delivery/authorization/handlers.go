package authorization

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

// ключ - сессия, значение - идентификатор пользователя
var sessionTable = make(map[uuid.UUID]uint64)

// ключ - идентификатор пользователя, значение - данные пользователя (экземпляр структуры)
var users = map[uint64]entity.User{}

var (
	sessionTableMu sync.RWMutex
	usersMu        sync.RWMutex
)

// устанавливает значение контекста user
func setUserContext(r *http.Request, user any) *http.Request {
	var ctx context.Context
	ctx = context.WithValue(r.Context(), "user", user)
	return r.WithContext(ctx)
}

func SessionAuthentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if loggedIn := !errors.Is(errNoSessionCookie, http.ErrNoCookie); loggedIn {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId != nil {
				// пользователь попортил куки | невозможно преобразовать в формат UUID
				r = setUserContext(r, nil)
				// !!! по идее нужно чистить куки пользователю !!!
			} else {
				// проверка на наличие UUID в таблице сессий
				sessionTableMu.RLock()
				userId, sessionExist := sessionTable[sessionId]
				sessionTableMu.RUnlock()

				if !sessionExist {
					r = setUserContext(r, nil)
				} else {
					usersMu.RLock()
					user := users[userId]
					usersMu.RUnlock()
					r = setUserContext(r, user)
				}
			}
		} else {
			// nil эквивалентно AnonymousUser
			r = setUserContext(r, nil)
		}
		handler.ServeHTTP(w, r)
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	return w
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	return w
}

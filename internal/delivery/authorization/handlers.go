package authorization

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

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
	// если пришел авторизованный пользователь, возвращаем 401
	user := r.Context().Value("user")
	if user != nil {
		http.Error(w, "Не хватает действительных учётных данных для целевого ресурса", http.StatusUnauthorized)
		return
	}

	// пользователь авторизовывается с помощью почты
	body, errWrongData := io.ReadAll(r.Body)
	if errWrongData != nil {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	// нужно декодироват данные, проверить credentionals и уже потом выдавать cookie
	email := body
	// нужно сделать еще проверку валидности введенных данных
	session_id := uuid.NewV4()
	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session_id.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	sessionTableMu.RLock()
	sessionTableMu[]
	sessionTableMu.RUnlock()
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	return w
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	return w
}

package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

type AuthHandler struct {
	sessionTableMu sync.RWMutex
	usersMu        sync.RWMutex
}

func (state *AuthHandler) SessionAuthentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if loggedIn := !errors.Is(errNoSessionCookie, http.ErrNoCookie); loggedIn {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId == nil {
				// проверка на наличие UUID в таблице сессий
				state.sessionTableMu.RLock()
				userEmail, sessionExist := state.sessionTable[sessionId]
				state.sessionTableMu.RUnlock()

				if sessionExist {
					state.usersMu.RLock()
					user := state.users[userEmail]
					state.usersMu.RUnlock()

					var ctx context.Context
					ctx = context.WithValue(r.Context(), "user", user)
					r.WithContext(ctx)
				}
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func (state *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	user := r.Context().Value("user")
	if user != nil {
		http.Error(w, "Не хватает действительных учётных данных для целевого ресурса", http.StatusUnauthorized)
		return
	}

	requestBody, errWrongData := io.ReadAll(r.Body)
	if errWrongData != nil {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	var bodyData entity.AuthorizationProps
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	if errRetrieveBodyData != nil {
		http.Error(w, "Ошибка при десериализации тела запроса", http.StatusBadRequest)
		return
	}

	if currentUser, userExist := state.users[bodyData.Email]; userExist && currentUser.CheckPassword(bodyData.Password) {
		sessionId := uuid.NewV4()
		// собираем Cookie
		expiration := time.Now().Add(14 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionId.String(),
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
		state.sessionTableMu.RLock()
		state.sessionTable[sessionId] = currentUser.Email
		state.sessionTableMu.RUnlock()

		// Собираем ответ
		w.Header().Set("Content-Type", "application/json")
		responseBody := map[string]string{"detail": "Пользователь успешно авторизован"}
		jsonResponse, err := json.Marshal(responseBody)
		if err != nil {
			http.Error(w, "Ошибка при сериализации тела ответа", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, errWriteResponseBody := w.Write(jsonResponse)
		if errWriteResponseBody != nil {
			http.Error(w, "Ошибка при формировании тела ответа", http.StatusBadRequest)
			return
		}
		return
	}
	http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
	return
}

func (state *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	// добавить валидацию на имя и пароль
	email := r.Form["email"][0]
	name := r.Form["name"][0]
	password := r.Form["password"][0]

	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9]{8}$`)
	if !regexPassword.MatchString(password) {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	regexName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{1,19}$`)
	if !regexName.MatchString(name) {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	regexEmail := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if regexEmail.MatchString(email) {
		state.users[email] = entity.User{Id: len(state.users), Email: email, Password: entity.HashData(password), Name: name}
	} else {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	sessionId := uuid.NewV4()
	state.sessionTable[sessionId] = email

	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	body, err := json.Marshal(state.users[email])
	if err != nil {
		http.Error(w, "Ошибка при сериализации тела ответа", http.StatusBadRequest)
		return
	}

	_, err = w.Write(body)

	if err != nil {
		http.Error(w, "Ошибка при формировании тела ответа", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (state *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	// если пришел неавторизованный пользователь, возвращаем 401
	user := r.Context().Value("user")
	if user == nil {
		http.Error(w, "Не хватает действительных учётных данных для целевого ресурса", http.StatusUnauthorized)
		return
	}

	// удаляем запись из таблицы сессий
	sessionCookie, errNoSessionCookie := r.Cookie("session_id")
	if errors.Is(errNoSessionCookie, http.ErrNoCookie) {
		http.Error(w, "Не хватает действительных учётных данных для целевого ресурса", http.StatusUnauthorized)
		return
	}
	// проверка на корректность UUID
	sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
	if errWrongSessionId != nil {
		http.Error(w, "Ошибка при получении ключа сессии", http.StatusBadRequest)
		return
	}
	state.sessionTableMu.RLock()
	delete(state.sessionTable, sessionId)
	state.sessionTableMu.RUnlock()

	// ставим заголовок для удаления сессионной куки в браузере
	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)
}

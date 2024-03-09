package authorization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/satori/uuid"

	"2024_1_kayros/internal/entity"
)

type AuthStore struct {
	SessionTable   map[uuid.UUID]string    // ключ - сессия, значение - идентификатор пользователя
	Users          map[string]*entity.User // ключ - почта пользователя, значение - данные пользователя (экземпляр структуры)
	SessionTableMu sync.RWMutex
	UsersMu        sync.RWMutex
}

type Registration struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewAuthStore() *AuthStore {
	users := []*entity.User{
		{Id: 1, Name: "Ivan", Email: "ivan@yandex.ru", Password: "358100c210df061db1f9a7a8945fa3140e169ddf67f7005c57c007647753e100"},
		{Id: 2, Name: "Sofia", Email: "sofia@yandex.ru"},
		{Id: 3, Name: "Bogdan", Email: "bogdan@yandex.ru"},
		{Id: 4, Name: "Pasha", Email: "pasha@yandex.ru"},
		{Id: 5, Name: "Ilya", Email: "ilya@yandex.ru"},
	}
	tmpUsers := map[string]*entity.User{}
	for _, user := range users {
		tmpUsers[user.Email] = user
	}
	return &AuthStore{
		SessionTable:   map[uuid.UUID]string{},
		Users:          tmpUsers,
		SessionTableMu: sync.RWMutex{},
		UsersMu:        sync.RWMutex{},
	}
}

func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		handler.ServeHTTP(w, r)
	})
}

func (state *AuthStore) SessionAuthentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, errNoSessionCookie := r.Cookie("session_id")
		if loggedIn := !errors.Is(errNoSessionCookie, http.ErrNoCookie); loggedIn {
			// проверка на корректность UUID
			sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
			if errWrongSessionId == nil {
				// проверка на наличие UUID в таблице сессий
				state.SessionTableMu.RLock()
				userEmail, sessionExist := state.SessionTable[sessionId]
				state.SessionTableMu.RUnlock()

				if sessionExist {
					state.UsersMu.RLock()
					user := state.Users[userEmail]
					state.UsersMu.RUnlock()

					var ctx context.Context
					ctx = context.WithValue(r.Context(), "user", user)
					r = r.WithContext(ctx)
				}
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func (state *AuthStore) SignIn(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	user := r.Context().Value("user")
	if user != nil {
		fmt.Print("ответ тело:", "aflafhas[vha[ishviashv[oasho[asovhas[ovhao[svha[osjvausviashv[uashvo[uahsobv123124124124124")
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
	_ = r.Body.Close()
	if errRetrieveBodyData != nil {
		http.Error(w, "Ошибка при десериализации тела запроса", http.StatusBadRequest)
		return
	}

	if currentUser, userExist := state.Users[bodyData.Email]; userExist && currentUser.CheckPassword(bodyData.Password) {
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

		state.SessionTableMu.RLock()
		state.SessionTable[sessionId] = currentUser.Email
		state.SessionTableMu.RUnlock()

		// Собираем ответ
		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(currentUser)
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

func (state *AuthStore) SignUp(w http.ResponseWriter, r *http.Request) {
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

	var bodyData Registration
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	_ = r.Body.Close()
	if errRetrieveBodyData != nil {
		http.Error(w, "Ошибка при десериализации тела запроса", http.StatusBadRequest)
		return
	}

	_, userAlreadyExist := state.Users[bodyData.Email]
	if userAlreadyExist {
		http.Error(w, "Пользователь с таким именем уже зарегистрирован", http.StatusBadRequest)
		return
	}

	regexPassword := regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	if !regexPassword.MatchString(bodyData.Password) {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	regexName := regexp.MustCompile(`^[a-zA-Zа-яА-Я][a-zA-Zа-яА-Я0-9]{1,19}$`)
	if !regexName.MatchString(bodyData.Name) {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	regexEmail := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if regexEmail.MatchString(bodyData.Email) {
		state.Users[bodyData.Email] = &entity.User{Id: len(state.Users), Email: bodyData.Email, Password: entity.HashData(bodyData.Password), Name: bodyData.Name}
	} else {
		http.Error(w, "Предоставлены неверные учетные данные", http.StatusBadRequest)
		return
	}

	sessionId := uuid.NewV4()
	state.SessionTable[sessionId] = bodyData.Email

	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	body, err := json.Marshal(state.Users[bodyData.Email])
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

func (state *AuthStore) SignOut(w http.ResponseWriter, r *http.Request) {
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
	state.SessionTableMu.RLock()
	delete(state.SessionTable, sessionId)
	state.SessionTableMu.RUnlock()

	// ставим заголовок для удаления сессионной куки в браузере
	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	// Успешно вышли из системы, возвращаем статус 200 OK и сообщение
	w.WriteHeader(http.StatusOK)
	message := "Пользователь успешно завершил сессию"
	_, errorWrite := w.Write([]byte(message))
	if errorWrite != nil {
		// Обработка ошибки записи сообщения в тело ответа
		http.Error(w, "Ошибка при формировании тела ответа", http.StatusBadRequest)
		return
	}
}

func (state *AuthStore) UserData(w http.ResponseWriter, r *http.Request) {
	// если пришел неавторизованный пользователь, возвращаем 401
	userPrt := r.Context().Value("user")
	if userPrt == nil {
		w = entity.ErrorResponse(w, "Не хватает прав для доступа", http.StatusBadRequest)
		return
	}
	user := userPrt.(*entity.User)
	response := entity.UserResponse{
		Id:   user.Id,
		Name: user.Name,
	}
	data, errSerialization := json.Marshal(response)
	if errSerialization != nil {
		w = entity.ErrorResponse(w, "Ошибка при формировании тела ответа", http.StatusBadRequest)
		return
	}
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

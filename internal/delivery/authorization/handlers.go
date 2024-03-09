package middlewares

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/satori/uuid"

	"2024_1_kayros/internal/entity"
)

type AuthHandler struct {
	Database entity.SystemDatabase
}

func (s *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
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

func (s *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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

func (s *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
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

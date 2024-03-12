package authorization

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
	db entity.SystemDatabase
}

func (state *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	authKey := r.Context().Value("authKey")
	if authKey != nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}

	requestBody, errWrongData := io.ReadAll(r.Body)
	if errWrongData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	var bodyData entity.AuthorizationProps
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	_ = r.Body.Close()
	if errRetrieveBodyData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	currentUser, userNotExist := state.db.Users.GetUser(bodyData.Email)
	if userNotExist == nil && currentUser.CheckPassword(bodyData.Password) {
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
		returnUser := state.Users[bodyData.Email]
		response := entity.UserResponse{
			Id:   returnUser.Id,
			Name: returnUser.Name,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

		_, errWriteResponseBody := w.Write(jsonResponse)
		if errWriteResponseBody != nil {
			w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
			return
		}
		return
	}
	w = entity.ErrorResponse(w, entity.BadAuthCredentials, http.StatusBadRequest)
	return
}

func isValidPassword(password string) bool {
	// Проверка на минимальную длину
	if len(password) < 8 {
		return false
	}

	// Проверка на наличие хотя бы одной буквы
	letterRegex := regexp.MustCompile(`[A-Za-z]`)
	if !letterRegex.MatchString(password) {
		return false
	}

	// Проверка на наличие хотя бы одной цифры
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password) {
		return false
	}

	// Проверка на наличие разрешенных символов
	validCharsRegex := regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	if !validCharsRegex.MatchString(password) {
		return false
	}

	return true
}

func (state *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	w.Header().Set("Content-Type", "application/json")
	user := r.Context().Value("user")
	if user != nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}

	requestBody, errWrongData := io.ReadAll(r.Body)
	if errWrongData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	var bodyData Registration
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	_ = r.Body.Close()
	if errRetrieveBodyData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	_, userAlreadyExist := state.Users[bodyData.Email]
	if userAlreadyExist {
		w = entity.ErrorResponse(w, entity.UserAlreadyExist, http.StatusBadRequest)
		return
	}

	if !isValidPassword(bodyData.Password) {
		w = entity.ErrorResponse(w, entity.BadRegCredentials, http.StatusBadRequest)
		return
	}

	regexName := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ][a-zA-Zа-яА-ЯёЁ0-9]{1,19}$`)
	if !regexName.MatchString(bodyData.Name) {
		w = entity.ErrorResponse(w, entity.BadRegCredentials, http.StatusBadRequest)
		return
	}

	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if regexEmail.MatchString(bodyData.Email) {
		state.Users[bodyData.Email] = &entity.User{
			Id:       len(state.Users),
			Email:    bodyData.Email,
			Password: entity.HashData(bodyData.Password),
			Name:     bodyData.Name,
		}
	} else {
		w = entity.ErrorResponse(w, entity.BadRegCredentials, http.StatusBadRequest)
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

	returnUser := state.Users[bodyData.Email]
	response := entity.UserResponse{
		Id:   returnUser.Id,
		Name: returnUser.Name,
	}
	body, err := json.Marshal(response)
	if err != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	_, err = w.Write(body)

	if err != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (state *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	// если пришел неавторизованный пользователь, возвращаем 401
	w.Header().Set("Content-Type", "application/json")
	user := r.Context().Value("user")
	if user == nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}

	// удаляем запись из таблицы сессий
	sessionCookie, errNoSessionCookie := r.Cookie("session_id")
	if errors.Is(errNoSessionCookie, http.ErrNoCookie) {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}
	// проверка на корректность UUID
	sessionId, errWrongSessionId := uuid.FromString(sessionCookie.Value)
	if errWrongSessionId != nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
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
	w = entity.ErrorResponse(w, "Сессия успешно завершена", http.StatusOK)
}

func (state *AuthHandler) UserData(w http.ResponseWriter, r *http.Request) {
	// если пришел неавторизованный пользователь, возвращаем 401
	w.Header().Set("Content-Type", "application/json")
	userPrt := r.Context().Value("user")
	if userPrt == nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}
	user := userPrt.(*entity.User)
	response := entity.UserResponse{
		Id:   user.Id,
		Name: user.Name,
	}
	data, errSerialization := json.Marshal(response)
	if errSerialization != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}
	_, errWrite := w.Write(data)
	if errWrite != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

package signin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

type AuthUsecase interface {
}

// должна возвращать интерфейс блока авторизации
func NewAuthBlock() {

}

func (state *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// если пришел авторизованный пользователь, возвращаем 401
	authKey := r.Context().Value("authKey")
	if authKey != nil {
		log.Println(entity.BadPermission)
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}

	requestBody, errWrongData := io.ReadAll(r.Body)
	if errWrongData != nil {
		log.Println(entity.UnexpectedServerError)
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusBadRequest)
		return
	}

	var bodyData entity.AuthorizationProps
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	_ = r.Body.Close()
	if errRetrieveBodyData != nil {
		log.Println(entity.BadRegCredentials)
		w = entity.ErrorResponse(w, entity.BadRegCredentials, http.StatusBadRequest)
		return
	}

	currentUser, userNotExist := state.DB.Users.GetUser(bodyData.Email)
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

		state.DB.Sessions.SetNewSession(sessionId, bodyData.Email)

		// Собираем ответ
		response := entity.UserResponse{
			Id:   currentUser.Id,
			Name: currentUser.Name,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println(entity.UnexpectedServerError)
			w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		_, errWriteResponseBody := w.Write(jsonResponse)
		if errWriteResponseBody != nil {
			log.Println(entity.UnexpectedServerError)
			w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusInternalServerError)
			return
		}
	} else {
		w = entity.ErrorResponse(w, entity.BadAuthCredentials, http.StatusBadRequest)
	}
}

func (state *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	w.Header().Set("Content-Type", "application/json")
	authKey := r.Context().Value("authKey")
	if authKey != nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}

	requestBody, errWrongData := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if errWrongData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusInternalServerError)
		return
	}

	var bodyData entity.RegistrationProps
	errRetrieveBodyData := json.Unmarshal(requestBody, &bodyData)
	if errRetrieveBodyData != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusInternalServerError)
		return
	}

	_, userNotExist := state.DB.Users.GetUser(bodyData.Email)
	if userNotExist == nil {
		w = entity.ErrorResponse(w, userNotExist.Error(), http.StatusBadRequest)
		return
	}

	if !entity.IsValidPassword(bodyData.Password) {
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
		hashedPassword, errHash := entity.HashData(bodyData.Password)
		if errHash != nil {
			w = entity.ErrorResponse(w, errHash.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = state.DB.Users.SetNewUser(bodyData.Email, entity.User{
			Id:       len(state.DB.Users.Data),
			Email:    bodyData.Email,
			Password: hashedPassword,
			Name:     bodyData.Name,
		})
	} else {
		w = entity.ErrorResponse(w, entity.BadRegCredentials, http.StatusBadRequest)
		return
	}

	sessionId := uuid.NewV4()
	state.DB.Sessions.SetNewSession(sessionId, bodyData.Email)

	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	returnUser, errGetUser := state.DB.Users.GetUser(bodyData.Email)
	if errGetUser != nil {
		w = entity.ErrorResponse(w, entity.UnexpectedServerError, http.StatusInternalServerError)
		return
	}
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
	authKey := r.Context().Value("authKey")
	fmt.Print(authKey)
	if authKey == nil {
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

	state.DB.Sessions.DeleteSession(sessionId)

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
	authKey := r.Context().Value("authKey")
	if authKey == nil {
		w = entity.ErrorResponse(w, entity.BadPermission, http.StatusUnauthorized)
		return
	}
	user, errGetUser := state.DB.Users.GetUser(authKey.(string))
	if errGetUser != nil {
		w = entity.ErrorResponse(w, errGetUser.Error(), http.StatusUnauthorized)
	}
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

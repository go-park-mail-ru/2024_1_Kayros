package authorization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/satori/uuid"

	"2024_1_kayros/internal/entity"
)

type AuthHandler struct {
	DB entity.AuthDatabase
}

// SignIn godoc @Summary Авторизация
// @Description Авторизация пользователя в системе. После авторизации устанавливается кука session_id, с помощью которой он в дальнейшем получает доступ к ресурсу
// @Tags User
// @Produce  json
// @Param		email		body		string					true	"email"
// @Param		password	body		string					true	"password"
// @Success 200 {object} 	entity.UserResponce 			"Пользователь успешно авторизован"
// @Failure 400 {object} 	entity.BadRegCredentials 		"Пользователь предоставил неверные данные для входа в аккаунт"
// @Failure 401 {object} 	entity.BadPermission			"Не хватает прав для доступа"
// @Failure 500 {object} 	entity.UnexpectedServerError	"Произошла неожиданная ошибка"
// @Router /signin [post]
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

// SignUp godoc @Summary Регистрация
// @Description Регистрация пользователя в системе. После регистрации пользователю устанавливается кука session_id, с помощью которой он в дальнейшем получает доступ к ресурсу
// @Tags User
// @Produce  json
// @Param		email		body		string					true	"email"
// @Param		password	body		string					true	"password"
// @Param		name		body		string					true	"name"
// @Param		phone		body		string					true	"phone"
// @Success 200 {object} 	entity.UserResponce 			"Пользователь успешно зарегистрирован"
// @Failure 400 {object} 	entity.BadRegCredentials 		"Были переданы некорректные данные для регистрации"
// @Failure 401 {object} 	entity.BadPermission			"Не хватает прав для доступа"
// @Failure 500 {object} 	entity.UnexpectedServerError	"Произошла неожиданная ошибка"
// @Router /signup [post]
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

// SignOut godoc @Summary Деавторизация
// @Description Выход пользователя из системы. Удаляется кука session_id из браузера. Данные доступны пользователю в режиме чтения
// @Tags User
// @Produce  json
// @Success 200 {object} 	entity.UserResponce 			"Пользователь успешно деавторизован"
// @Failure 400 {object} 	entity.BadRegCredentials 		"Пользователь не может деавторизоваться"
// @Failure 401 {object} 	entity.BadPermission			"Неудачная попытка покинуть систему"
// @Failure 500 {object} 	entity.UnexpectedServerError	"Произошла неожиданная ошибка"
// @Router /signout [post]
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

// UserData godoc @Summary Данные пользователя
// @Description Получение данных о пользователе
// @Tags User
// @Produce  json
// @Success 200 {object} 	entity.UserResponce 			"Данные пользователь успешно отправлены"
// @Failure 401 {object} 	entity.BadPermission			"Неудачная получить данные"
// @Failure 500 {object} 	entity.UnexpectedServerError	"Произошла неожиданная ошибка"
// @Router /signout [post]
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

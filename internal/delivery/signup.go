package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"

	"2024_1_kayros/internal/entity"
	"github.com/satori/uuid"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
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

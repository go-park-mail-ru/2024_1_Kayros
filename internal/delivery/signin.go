package delivery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/satori/uuid"

	"2024_1_kayros/internal/entity"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
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

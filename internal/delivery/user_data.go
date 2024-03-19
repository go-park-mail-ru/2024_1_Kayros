package delivery

import (
	"encoding/json"
	"net/http"

	"2024_1_kayros/internal/entity"
)

func UserData(w http.ResponseWriter, r *http.Request) {
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

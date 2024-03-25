package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/session"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/satori/uuid"
)

type AuthUsecaseInterface interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	SignUpUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
	SignOut(w http.ResponseWriter, r *http.Request)
}

type AuthUsecase struct {
	repoUser    user.UserRepositoryInterface
	repoSession session.SessionRepositoryInterface
}

func NewAuthUsecase(repoUserProps user.UserRepositoryInterface, repoSessionProps session.SessionRepositoryInterface) AuthUsecaseInterface {
	return &AuthUsecase{
		repoUser:    repoUserProps,
		repoSession: repoSessionProps,
	}
}

func (state *AuthUsecase) SignIn(w http.ResponseWriter, r *http.Request) {

}

// пока что логгера нет, нужно будет пробрасывать
func (uc *AuthUsecase) SignUpUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	requestBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
		return w
	}

	// нужно на атрибуты DTO навесит теги валидации
	var bodyDataDTO dto.SignUpDTO
	err = json.Unmarshal(requestBody, &bodyDataDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	// нужно будет добавить метод на существование пользователя
	_, err = uc.repoUser.GetByEmail(bodyDataDTO.Email)
	if err == nil {
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
		return w
	}

	// пока что проверка валидации через вспомогательную функцию
	if !functions.IsValidPassword(bodyDataDTO.Password) {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	regexName := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ][a-zA-Zа-яА-ЯёЁ0-9]{1,19}$`)
	if !regexName.MatchString(bodyDataDTO.Name) {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if regexEmail.MatchString(bodyDataDTO.Email) {
		u := &entity.User{
			Name:     bodyDataDTO.Name,
			Email:    bodyDataDTO.Email,
			Password: bodyDataDTO.Password,
			Phone:    bodyDataDTO.Phone,
		}

		u.Password, err = functions.HashData(u.Password)
		if err != nil {
			fmt.Println(err)
			w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
			return w
		}

		// тут стоит выводить текст возвращаемой ошибки и возвращать другую ошибку в handler
		_, err = uc.repoUser.Create(u)
		if err != nil {
			w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
			return w
		}
	} else {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	sessionId := uuid.NewV4().String()
	err = uc.repoSession.SetValue(alias.SessionKey(sessionId), alias.SessionValue(bodyDataDTO.Email))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
		return w
	}

	// собираем Cookie
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	userByEmail, err := uc.repoUser.GetByEmail(bodyDataDTO.Email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
		return w
	}

	// нужно:
	// 1) написать DTO структуры для маршала
	// 2) присвоить им userByEmail (написать функцию, которая перенесет запись из entity в dto
	body, err := json.Marshal(userByEmail)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.IntServerError, http.StatusInternalServerError)
		return w
	}
	w.WriteHeader(http.StatusOK)
	return w
}

func (state *AuthUsecase) SignOut(w http.ResponseWriter, r *http.Request) {
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
		w = entity.ErrorResponse(w, entity.InternalServerError, http.StatusBadRequest)
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

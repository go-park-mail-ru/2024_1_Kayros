package auth

import (
	"encoding/json"
	"io"
	"net/http"
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

type Usecase interface {
	SignInUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
	SignUpUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
	SignOutUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter
}

type UsecaseLayer struct {
	repoUser    user.UserRepositoryInterface
	repoSession session.SessionRepositoryInterface
}

func NewUsecase(repoUserProps user.UserRepositoryInterface, repoSessionProps session.SessionRepositoryInterface) Usecase {
	return &UsecaseLayer{
		repoUser:    repoUserProps,
		repoSession: repoSessionProps,
	}
}

func (uc *UsecaseLayer) SignInUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	// если пришел авторизованный пользователь, возвращаем 401
	userEmail := r.Context().Value("email")
	if userEmail != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	var bodyDTO dto.SignInDTO
	err = json.Unmarshal(body, &bodyDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	isValid := bodyDTO.Validate()
	if !isValid {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	u, err := uc.repoUser.GetByEmail(bodyDTO.Email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
		return w
	}

	isRightPassword, err := uc.repoUser.CheckPassword(alias.UserId(u.Id), bodyDTO.Password)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	if isRightPassword {
		sessionId := uuid.NewV4()
		expiration := time.Now().Add(14 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionId.String(),
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)

		err = uc.repoSession.SetValue(alias.SessionKey(sessionId.String()), alias.SessionValue(u.Email))
		if err != nil {
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return w
		}

		// Собираем ответ
		uDTO := &dto.UserDTO{
			Id:       u.Id,
			Name:     u.Name,
			Phone:    u.Phone,
			Email:    u.Email,
			ImgUrl:   u.ImgUrl,
			Password: u.Password,
		}
		jsonResponse, err := json.Marshal(uDTO)
		if err != nil {
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return w
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return w
		}
	} else {
		w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
		return w
	}
	return w
}

// SignUpUser пока что логгера нет, нужно будет пробрасывать
func (uc *UsecaseLayer) SignUpUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	authKey := r.Context().Value("authKey")
	if authKey != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	requestBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	var bodyDataDTO dto.SignUpDTO
	err = json.Unmarshal(requestBody, &bodyDataDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	isValid := bodyDataDTO.Validate()
	if !isValid {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return w
	}

	isExist := uc.repoUser.IsExistByEmail(bodyDataDTO.Email)
	if isExist {
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
		return w
	}

	u := &entity.User{
		Name:     bodyDataDTO.Name,
		Email:    bodyDataDTO.Email,
		Password: bodyDataDTO.Password,
	}

	err = uc.repoUser.Create(u)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	sessionId := uuid.NewV4().String()
	err = uc.repoSession.SetValue(alias.SessionKey(sessionId), alias.SessionValue(bodyDataDTO.Email))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
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

	u, err = uc.repoUser.GetByEmail(bodyDataDTO.Email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	uDTO := &dto.UserDTO{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		ImgUrl:   u.ImgUrl,
		Password: u.Password,
	}
	body, err := json.Marshal(uDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}
	w.WriteHeader(http.StatusOK)
	return w
}

func (uc *UsecaseLayer) SignOutUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	authKey := r.Context().Value("authKey")
	if authKey == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	// проверка на соответствие формату UUID4
	sessionKey := sessionCookie.Value
	_, err = uuid.FromString(sessionKey)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	err = uc.repoSession.DeleteKey(alias.SessionKey(sessionKey))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	// ставим заголовок для удаления сессионной куки в браузере
	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	// Успешно вышли из системы, возвращаем статус 200 OK и сообщение
	w.WriteHeader(http.StatusOK)
	w = functions.ErrorResponse(w, "Сессия успешно завершена", http.StatusOK)
	return w
}

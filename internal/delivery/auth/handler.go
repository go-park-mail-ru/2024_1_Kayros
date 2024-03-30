package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/satori/uuid"
)

type Delivery struct {
	ucSession session.Usecase
	ucUser    user.Usecase
}

func NewDeliveryLayer(sessionUsecase session.Usecase, userUsecase user.Usecase) *Delivery {
	return &Delivery{
		ucSession: sessionUsecase,
		ucUser:    userUsecase,
	}
}

// SignUpUser пока что логгера нет, нужно будет пробрасывать
func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email")
	if email != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var bodyDataDTO dto.SignUpDTO
	err = json.Unmarshal(requestBody, &bodyDataDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid := bodyDataDTO.Validate()
	if !isValid {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isExist, err := d.ucUser.IsExistByEmail(r.Context(), bodyDataDTO.Email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isExist {
		w = functions.ErrorResponse(w, myerrors.UserAlreadyExistError, http.StatusBadRequest)
		return
	}

	u := &entity.User{
		Name:     bodyDataDTO.Name,
		Email:    bodyDataDTO.Email,
		Password: bodyDataDTO.Password,
	}

	// получили данные пользователя из БД
	u, err = d.ucUser.Create(r.Context(), u)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	sessionId := uuid.NewV4().String()
	_, err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId), alias.SessionValue(bodyDataDTO.Email))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
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

	uDTO := functions.ConvIntoUserDTO(u)
	w = functions.JsonResponse(w, uDTO)
	w.WriteHeader(http.StatusOK)
	return
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	// если пришел авторизованный пользователь, возвращаем 401
	email := r.Context().Value("email")
	if email != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	var bodyDTO dto.SignInDTO
	err = json.Unmarshal(body, &bodyDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	isValid := bodyDTO.Validate()
	if !isValid {
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}

	u, err := d.ucUser.GetByEmail(r.Context(), bodyDTO.Email)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
		return
	}

	isEqual, err := d.ucUser.CheckPassword(r.Context(), bodyDTO.Email, bodyDTO.Password)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}

	if isEqual {
		sessionId := uuid.NewV4()
		expiration := time.Now().Add(14 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionId.String(),
			Expires:  expiration,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)

		_, err = d.ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(u.Email))
		if err != nil {
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}

		// Собираем ответ
		uDTO := functions.ConvIntoUserDTO(u)
		w = functions.JsonResponse(w, uDTO)
		w.WriteHeader(http.StatusOK)
		return
	}
	w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsError, http.StatusBadRequest)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	authKey := r.Context().Value("email")
	if authKey == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	_, err = d.ucSession.DeleteKey(r.Context(), alias.SessionKey(sessionCookie.Value))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
	}

	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	w = functions.JsonResponse(w, "Сессия успешно завершена")
	w.WriteHeader(http.StatusOK)
}

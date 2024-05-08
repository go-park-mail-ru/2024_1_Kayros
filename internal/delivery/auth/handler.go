package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/auth"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/props"
	"2024_1_kayros/internal/utils/sanitizer"
	authv1 "2024_1_kayros/microservices/auth/proto"

	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase
	ucUser    user.Usecase
	ucCsrf    session.Usecase
	ucAuth    auth.Usecase
	logger    *zap.Logger
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, ucAuthProps auth.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucUser:    ucUserProps,
		ucCsrf:    ucCsrfProps,
		ucAuth:    ucAuthProps,
		logger:    loggerProps,
		cfg:       cfgProps,
	}
}

func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	unauthId := functions.GetCtxUnauthId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.RegisteredRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var signupDTO dto.UserSignUp
	err = json.Unmarshal(body, &signupDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := signupDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	u := dto.NewUserFromSignUpForm(&signupDTO)

	signUpData := &authv1.SignUpCredentials {
		Email: u.Email,
		UnauthId: unauthId,
		Name: u.Name,
		Password: u.Password,
	}
	uAuth, err := d.ucAuth.SignUp(r.Context(), signUpData)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.UserAlreadyExist) {
			w = functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		// error `myerrors.SqlNoRowsUserRelation` is handled
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	u = sanitizer.User(cnvAuthUserIntoEntityUser(uAuth))
	uDTO := dto.NewUserData(u)

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, u.Email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func cnvUserIntoAuthUser (u *entity.User) *authv1.User {
	return &authv1.User{
		Id: u.Id,
		Name: u.Name,
		Phone: u.Phone,
		Email: u.Email,
		Address: u.Address,
		ImgUrl: u.ImgUrl,
		CardNumber: u.CardNumber,
		Password: u.Password,
	}
}


func cnvAuthUserIntoEntityUser (u *authv1.User) *entity.User {
	return &entity.User{
		Id: u.GetId(),
		Name: u.GetName(),
		Phone: u.GetPhone(),
		Email: u.GetEmail(),
		Address: u.GetAddress(),
		ImgUrl: u.GetImgUrl(),
		CardNumber: u.GetCardNumber(),
		Password: u.GetPassword(),
	}
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	unauthId := functions.GetCtxUnauthId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var bodyDTO dto.UserSignIn
	err = json.Unmarshal(body, &bodyDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := bodyDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}


	signInData := &authv1.SignInCredentials {
		Email: bodyDTO.Email,
		UnauthId: unauthId,
		Password: bodyDTO.Password,
	}
	u, err := d.ucAuth.SignIn(r.Context(),signInData)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) || errors.Is(err, myerrors.BadAuthPassword) {
			w = functions.ErrorResponse(w, myerrors.BadAuthCredentialsRu, http.StatusBadRequest)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uSanitizer := sanitizer.User(cnvAuthUserIntoEntityUser(u))
	uDTO := dto.NewUserData(uSanitizer)

	setCookieProps := props.GetSetCookieProps(d.ucCsrf, d.ucSession, u.Email, d.cfg.CsrfSecretKey)
	w, err = functions.SetCookie(w, r, setCookieProps)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		w = functions.ErrorResponse(w, myerrors.SignOutAlreadyRu, http.StatusUnauthorized)
		return
	}

	w, err := functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Сессия успешно завершена"})
}
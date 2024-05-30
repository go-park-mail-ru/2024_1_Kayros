package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/auth"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/sanitizer"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Delivery struct {
	ucSession session.Usecase // methods for communicating to microservice session
	ucAuth    auth.Usecase
	logger    *zap.Logger
	ucUser    user.Usecase
	cfg       *config.Project
	metrics   *metrics.Metrics
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucAuthProps auth.Usecase, ucUserProps user.Usecase, loggerProps *zap.Logger, metrics *metrics.Metrics) *Delivery {
	return &Delivery{
		ucSession: ucSessionProps,
		ucAuth:    ucAuthProps,
		ucUser:    ucUserProps,
		logger:    loggerProps,
		cfg:       cfgProps,
		metrics:   metrics,
	}
}

func (d *Delivery) SignUp(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.RegisteredRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var signupDTO dto.UserSignUp
	err = easyjson.Unmarshal(body, &signupDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := signupDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	u := dto.NewUserFromSignUpForm(&signupDTO)
	uSignedUp, err := d.ucAuth.SignUp(r.Context(), u)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.UserAlreadyExist) {
			functions.ErrorResponse(w, myerrors.UserAlreadyExistRu, http.StatusBadRequest)
			return
		}
		// error `myerrors.SqlNoRowsUserRelation` is handled
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uDTO := dto.NewUserData(sanitizer.User(uSignedUp))

	w, err = functions.SetCookie(w, r, d.ucSession, u.Email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignIn(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.AuthorizedRu, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	var bodyDTO dto.UserSignIn
	err = easyjson.Unmarshal(body, &bodyDTO)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	isValid, err := bodyDTO.Validate()
	if err != nil || !isValid {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	u, err := d.ucAuth.SignIn(r.Context(), bodyDTO.Email, bodyDTO.Password)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) || errors.Is(err, myerrors.BadAuthPassword) {
			functions.ErrorResponse(w, myerrors.BadAuthCredentialsRu, http.StatusBadRequest)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	uDTO := dto.NewUserData(sanitizer.User(u))

	w, err = functions.SetCookie(w, r, d.ucSession, u.Email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	functions.JsonResponse(w, uDTO)
}

func (d *Delivery) SignOut(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		functions.ErrorResponse(w, myerrors.SignOutAlreadyRu, http.StatusUnauthorized)
		return
	}

	w, err := functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Сессия успешно завершена"})
}


func (d *Delivery) AuthVk(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)

	requestBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, errors.New("Invalid payload"), http.StatusBadRequest)
		return
	}

	var data dto.VkBodyRequest
	err = easyjson.Unmarshal(requestBody, &data)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, errors.New("Invalid JSON in payload") , http.StatusBadRequest)
		return
	}

	uuid := data.Payload.UUID
	silentToken := data.Payload.Token
	if uuid == "" || silentToken == "" {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, errors.New("Missing uuid or token in payload") , http.StatusBadRequest)
		return
	}

    vkURL := fmt.Sprintf("https://api.vk.com/method/auth.exchangeSilentAuthToken?v=5.131&token=%s&access_token=%s&uuid=%s", silentToken, d.cfg.Oauth.AccessToken, uuid)

    resp, err := http.Get(vkURL)
    if err != nil {
        d.logger.Error("VK API request failed", zap.Error(err))
        functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        d.logger.Error("Failed to read VK API response", zap.Error(err))
        functions.ErrorResponse(w, errors.New("Failed to read VK API response"), http.StatusBadRequest)
        return
    }

    var vkResponse map[string]interface{}
    err = json.Unmarshal(responseBody, &vkResponse)
    if err != nil {
        d.logger.Error("Failed to parse VK API response", zap.Error(err))
        functions.ErrorResponse(w, errors.New("Failed to parse VK API response"), http.StatusBadRequest)
        return
    }

	fmt.Println()
	fmt.Println()
	fmt.Println()
	d.logger.Info(fmt.Sprintf("%v", vkResponse))
	re, ok := vkResponse["response"].(map[string]interface{})
	if !ok {
		d.logger.Error("Failed to parse VK API response", zap.Error(err))
        functions.ErrorResponse(w, errors.New("Failed to 1 VK API response data"), http.StatusBadRequest)
        return
	}
	email, ok := re["email"].(string)
	if !ok {
		d.logger.Error("Failed to parse VK API response", zap.Error(err))
        functions.ErrorResponse(w, errors.New("Failed to retrieve email from VK API response"), http.StatusBadRequest)
        return
	}

	var userDto *dto.UserGet
	_, err = d.ucUser.GetData(r.Context(), email) 
	if err != nil {
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		fmt.Println(email)
		fmt.Println(data.Payload.User.FirstName + data.Payload.User.LastName)
		fmt.Println(data.Payload.User.Avatar)
		userDB, err := d.ucAuth.SignUp(r.Context(), &entity.User{
			Email: email,
			Name: data.Payload.User.FirstName + data.Payload.User.LastName,
			Password: "",
			ImgUrl: data.Payload.User.Avatar,
		})
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		userDto = dto.NewUserData(userDB)
	} else {
		userDB, err := d.ucAuth.SignIn(r.Context(), email, "")
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		userDto = dto.NewUserData(userDB)
	}

	w, err = functions.SetCookie(w, r, d.ucSession, email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}
	functions.JsonResponse(w, userDto)
}

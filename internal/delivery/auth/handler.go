package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	payload := r.URL.Query().Get("payload")

	payloadURL, err := url.QueryUnescape(payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

    // Используем стандартную библиотеку для работы с JSON, например, encoding/json
	var data map[string]interface{}
	err = json.Unmarshal([]byte(payloadURL), &data)
	if err != nil {
		http.Error(w, "Invalid JSON in payload", http.StatusBadRequest)
		return
	}

	 // Извлекаем данные пользователя
	user, ok := data["user"].(map[string]interface{})
	if !ok {
		fmt.Println("User data not found")
		return
	}

	// Извлекаем имя и фамилию
	firstName, ok := user["first_name"].(string)
	if !ok {
		fmt.Println("First name not found")
		return
	}
   
	lastName, ok := user["last_name"].(string)
	if !ok {
		fmt.Println("Last name not found")
		return
	}
 

	// Извлекаем uuid и token из распарсенного JSON
	uuid, ok1 := data["uuid"].(string)
	silentToken, ok2 := data["token"].(string)
	if !ok1 || !ok2 {
		http.Error(w, "Missing uuid or token in payload", http.StatusBadRequest)
		return
	}

    // VK API request
    accessToken := "43ab9f8b43ab9f8b43ab9f8bd340b3b4e4443ab43ab9f8b259a8b25328287013ece3d21"
    vkURL := fmt.Sprintf("https://api.vk.com/method/auth.exchangeSilentAuthToken?v=5.131&token=%s&access_token=%s&uuid=%s", silentToken, accessToken, uuid)

    resp, err := http.Get(vkURL)
    if err != nil {
        d.logger.Error("VK API request failed", zap.Error(err))
        functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read VK API response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        d.logger.Error("Failed to read VK API response", zap.Error(err))
        functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
        return
    }

    var vkResponse map[string]interface{}
    err = json.Unmarshal(body, &vkResponse)
    if err != nil {
        d.logger.Error("Failed to parse VK API response", zap.Error(err))
        functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
        return
    }

    // Log the VK API response
    d.logger.Info("VK API Response", zap.Any("vkResponse", vkResponse))

	var email string
	if response, ok := vkResponse["response"].(map[string]interface{}); ok {
        if email, ok = response["email"].(string); ok {
            fmt.Println("Email:", email)
        }
    }

	userDB, err := d.ucUser.GetData(r.Context(), email)
	if err != nil {
		d.ucAuth.SignUp(r.Context(), &entity.User{
			Email: email,
			Name: lastName + firstName,
			Password: "qwerty12345",
		})
	}

	w, err = functions.SetCookie(w, r, d.ucSession, email, d.cfg)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
	}

	userDto := dto.NewUserData(userDB)
	functions.JsonResponse(w, userDto)
}

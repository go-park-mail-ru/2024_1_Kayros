package user

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/usecase/user"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/sanitizer"
	"github.com/asaskevich/govalidator"
)

type Delivery struct {
	ucSession session.Usecase
	ucCsrf    session.Usecase
	ucUser    user.Usecase
	cfg       *config.Project
}

func NewDeliveryLayer(cfgProps *config.Project, ucSessionProps session.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase) *Delivery {
	return &Delivery{
		ucUser:    ucUserProps,
		ucSession: ucSessionProps,
		ucCsrf:    ucCsrfProps,
		cfg:       cfgProps,
	}
}

func (d *Delivery) UserData(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	u, err := d.ucUser.GetUserData(r.Context(), email, requestId, logger)
	if err != nil {
		// нужно обработать кастомные ошибки БД и тут их обрабатывать
		return
	}

	uSanitizer := sanitizer.User(u)
	uDTO := dto.NewUserData(uSanitizer)
	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	file, handler, u, err := dto.GetUpdatedUserData(r)
	defer func(file multipart.File) {
		if file != nil {
			err = file.Close()
			if err != nil {
				log.Println("Error of closing file")
			}
		}
	}(file)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			err = errors.New("Некорректный пароль")
			return
		}
		if strings.Contains(err.Error(), "phone") {
			err = errors.New("Некорректный номер телефона")
			return
		}
		if strings.Contains(err.Error(), "name") {
			err = errors.New("Некорректное имя")
			return
		}
		err = errors.New("Некорректные данные")
		return
	}

	uUpdated, err := d.ucUser.UpdateUserData(r.Context(), email, file, handler, u, requestId, logger)
	if err != nil {
		return
	}
	uSanitizer := sanitizer.User(uUpdated)
	userDTO := dto.NewUserData(uSanitizer)

	err = functions.DeleteCookies(r, d.ucSession, d.ucCsrf)
	if err != nil {
		return
	}
	functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)
	w = functions.JsonResponse(w, userDTO)
}

func (d *Delivery) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	var address addressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		return
	}
	isValid, err := address.Validate()
	if err != nil {
		return
	}
	if !isValid {
		return
	}
	uUpdated, err := d.ucUser.UpdateUserAddress(r.Context(), email, address.Data, requestId, logger)
	if err != nil {
		return
	}

	uSanitizer := sanitizer.User(uUpdated)
	uDTO := dto.NewUserData(uSanitizer)
	err = functions.DeleteCookies(r, d.ucSession, d.ucCsrf)
	if err != nil {
		return
	}
	functions.SetCookie(w, r, d.ucCsrf, d.ucSession, email, d.cfg.CsrfSecretKey)

	w = functions.JsonResponse(w, uDTO)
}

func (d *Delivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	logger, err := functions.GetCtxLogger(r)
	if err != nil {
		return
	}
	requestId, err := functions.GetCtxRequestId(r)
	if err != nil {
		return
	}
	email, err := functions.GetCtxEmail(r)
	if err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	var password passwordData
	err = json.Unmarshal(body, &password)
	if err != nil {
		return
	}

	isValid, err := password.Validate()
	if err != nil || !isValid {
		return
	}

	w = functions.JsonResponse(w, map[string]string{"detail": "Пароль был успешно обновлен"})
	functions.LogOkResponse(d.logger, requestId, cnst.NameHandlerUpdateUser, cnst.DeliveryLayer)
}

type addressData struct {
	Data string `json:"user_address_domain"`
}

func (d *addressData) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type passwordData struct {
	Data string `json:"password" valid:"user_pwd"`
}

func (d *passwordData) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

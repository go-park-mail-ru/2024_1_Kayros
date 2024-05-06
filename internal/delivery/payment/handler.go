package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/usecase/order"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type Payment struct {
	logger    *zap.Logger
	ucOrder   order.Usecase
	ucCsrf    session.Usecase
	ucSession session.Usecase
	cfg       *config.Payment
}

func NewPaymentDelivery(loggerProps *zap.Logger, ucOrderProps order.Usecase, ucCsrfProps session.Usecase, ucSessionProps session.Usecase, cfgProps *config.Payment) *Payment {
	return &Payment{
		logger:    loggerProps,
		ucOrder:   ucOrderProps,
		ucCsrf:    ucCsrfProps,
		ucSession: ucSessionProps,
		cfg:       cfgProps,
	}
}

func (d *Payment) OrderGetPayUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		d.logger.Error(myerrors.AuthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	basket, err := d.ucOrder.GetBasket(r.Context(), email)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			w, err = functions.FlashCookie(r, w, d.ucCsrf, d.ucSession)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, myerrors.SqlNoRowsOrderRelation) {
			w = functions.ErrorResponse(w, myerrors.NoBasketRu, http.StatusNotFound)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	// we need to retrieve basket sum
	bodyRequestYooMoney := fmt.Sprintf(`{
		"amount": {
		"value": %d,
			"currency": "RUB"
		},
		"payment_method_data": {
		"type": "bank_card",
			"card": {
			"cardholder": "MR CARDHOLDER",
				"csc": "213",
				"expiry_month": "12",
				"expiry_year": "2024",
				"number": "5555555555554477"
		}
		},
		"capture": true,
		"confirmation": {
		"type": "redirect",
			"return_url": "https://resto-go.ru"
		},
		"description": "Заказ №%d"
	}`, basket.Sum, basket.Id)
	requestBody := bytes.NewBuffer([]byte(bodyRequestYooMoney))

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", requestBody)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewV4().String())
	req.SetBasicAuth(d.cfg.StoreId, d.cfg.SecretKey)

	resp, err := client.Do(req)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	// receiving 'confirmation_url'
	confirmationURL := data["confirmation"].(map[string]interface{})["confirmation_url"].(string)
	w = functions.JsonResponse(w, map[string]string{"url": confirmationURL})
}

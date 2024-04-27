package statistic

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/statistic"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	ucQuiz statistic.Usecase
	ucUser user.Usecase
	logger *zap.Logger
}

func NewDeliveryLayer(ucQuizProps statistic.Usecase, ucUserProps user.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucQuiz: ucQuizProps,
		ucUser: ucUserProps,
		logger: loggerProps,
	}
}

func (d *Delivery) GetStatistic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctxRequestId := r.Context().Value("request_id")
	requestId := ""
	if ctxRequestId == nil {
		d.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email == "" {
		d.logger.Error(myerrors.UnauthorizedError, zap.String("request_id", requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	stats, err := d.ucQuiz.GetStatistic(r.Context())
	if err != nil {
		d.logger.Error(err.Error(), zap.String("request_id", requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return
	}

	w = functions.JsonResponse(w, stats)
}

func (d *Delivery) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctxRequestId := r.Context().Value("request_id")
	requestId := ""
	if ctxRequestId == nil {
		d.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}
	url := r.URL.Query().Get("url") //    /restaurants, /address
	qs := []*entity.Question{}
	var err error
	if url != "" {
		qs, err = d.ucQuiz.GetQuestionInfo(r.Context(), url)
		if err != nil {
			d.logger.Error(err.Error(), zap.String("request_id", requestId))
			w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
			return
		}
	}
	qDTO := dto.QuestionReturn(qs)
	w = functions.JsonResponse(w, qDTO)
}

//func (d *Delivery) UpdateAnswer(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	requestId := ""
//	ctxRequestId := r.Context().Value("request_id")
//	if ctxRequestId == nil {
//		d.logger.Error("request_id передан не был")
//	} else {
//		requestId = ctxRequestId.(string)
//	}
//
//	var qi dto.Question
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		d.logger.Error(err.Error(), zap.String("request_id", requestId))
//		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
//		return
//	}
//	if err = r.Body.Close(); err != nil {
//		d.logger.Error(err.Error(), zap.String("request_id", requestId))
//		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
//		return
//	}
//
//	err = json.Unmarshal(body, &qi)
//	if err != nil {
//		d.logger.Error(err.Error(), zap.String("request_id", requestId))
//		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
//		return
//	}
//
//	isValid, err := qi.Validate()
//	if !isValid || err != nil {
//		d.logger.Error(err.Error(), zap.String("request_id", requestId))
//		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
//		return
//	}
//
//	email := ""
//	ctxEmail := r.Context().Value("email")
//	if ctxEmail != nil {
//		email = ctxEmail.(string)
//	}
//	if email != "" {
//		u, err := d.ucUser.GetByEmail(r.Context(), email)
//		if err != nil {
//			d.logger.Error(err.Error(), zap.String("request_id", requestId))
//			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
//			return
//		}
//		err = d.ucQuiz.Update(r.Context(), qi.Id, qi.Rating, strconv.Itoa(int(u.Id)))
//	} else {
//		token := ""
//		ctxToken := r.Context().Value("unauth_token")
//		if ctxToken != nil {
//			token = ctxToken.(string)
//			err = d.ucQuiz.Update(r.Context(), qi.Id, qi.Rating, token)
//		}
//		if err != nil {
//			d.logger.Error(err.Error(), zap.String("request_id", requestId))
//			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
//			return
//		}
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

func (d *Delivery) AddAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := ""
	ctxRequestId := r.Context().Value("request_id")
	if ctxRequestId == nil {
		d.logger.Error("request_id передан не был")
	} else {
		requestId = ctxRequestId.(string)
	}

	var qi []*dto.QuestionInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error(err.Error(), zap.String("request_id", requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	if err = r.Body.Close(); err != nil {
		d.logger.Error(err.Error(), zap.String("request_id", requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &qi)
	if err != nil {
		d.logger.Error(err.Error(), zap.String("request_id", requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsError, http.StatusBadRequest)
		return
	}
	email := ""
	ctxEmail := r.Context().Value("email")
	if ctxEmail != nil {
		email = ctxEmail.(string)
	}
	if email != "" {
		u, err := d.ucUser.GetByEmail(r.Context(), email)
		if err != nil {
			d.logger.Error(err.Error(), zap.String("request_id", requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
			return
		}
		for _, q := range qi {
			err = d.ucQuiz.Create(r.Context(), q.Id, q.Rating, strconv.Itoa(int(u.Id)))
			if err != nil {
				d.logger.Error(err.Error(), zap.String("request_id", requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
				return
			}
		}
	} else {
		token := ""
		ctxToken := r.Context().Value("unauth_token")
		if ctxToken != nil {
			token = ctxToken.(string)
			for _, q := range qi {
				err = d.ucQuiz.Create(r.Context(), q.Id, q.Rating, token)
				if err != nil {
					d.logger.Error(err.Error(), zap.String("request_id", requestId))
					w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
					return
				}
			}
		}
	}
	w = functions.JsonResponse(w, map[string]string{"detail": "Пользователь успешно проголосовал"})
}

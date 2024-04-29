package statistic

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/usecase/statistic"
	"2024_1_kayros/internal/usecase/user"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	ucQuiz    statistic.Usecase
	ucUser    user.Usecase
	ucCsrf    session.Usecase
	ucSession session.Usecase
	logger    *zap.Logger
}

func NewDeliveryLayer(ucQuizProps statistic.Usecase, ucUserProps user.Usecase, ucCsrfProps session.Usecase, ucSessionProps session.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		ucQuiz:    ucQuizProps,
		ucUser:    ucUserProps,
		logger:    loggerProps,
		ucCsrf:    ucCsrfProps,
		ucSession: ucSessionProps,
	}
}

func (d *Delivery) GetStatistic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	stats, err := d.ucQuiz.GetStatistic(r.Context())
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusUnauthorized)
		return
	}

	w = functions.JsonResponse(w, stats)
}

func (d *Delivery) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)

	url := r.URL.Query().Get("url")
	qs := []*entity.Question{}
	var err error
	if url != "" {
		qs, err = d.ucQuiz.GetQuestionsOnFocus(r.Context(), url)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}
	qDTO := dto.QuestionReturn(qs)
	w = functions.JsonResponse(w, qDTO)
}

func (d *Delivery) AddAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)

	var qi []*dto.QuestionInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &qi)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	hasVoted := false
	email := functions.GetCtxEmail(r)
	unauthId := functions.GetCtxUnauthId(r)
	if email != "" {
		u, err := d.ucUser.GetData(r.Context(), email)
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
			w = functions.ErrorResponse(w, myerrors.InternalServerEn, http.StatusInternalServerError)
			return
		}
		for _, q := range qi {
			err = d.ucQuiz.Create(r.Context(), q.Id, q.Rating, strconv.Itoa(int(u.Id)))
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				if errors.Is(err, myerrors.QuizAdd) {
					w = functions.ErrorResponse(w, myerrors.QuizAddRu, http.StatusInternalServerError)
				} else {
					w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				}
				return
			}
		}
		hasVoted = true
	} else if unauthId != "" {
		for _, q := range qi {
			err = d.ucQuiz.Create(r.Context(), q.Id, q.Rating, unauthId)
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
		}
		hasVoted = true
	}
	if hasVoted {
		w = functions.JsonResponse(w, map[string]string{"detail": "Пользователь успешно проголосовал"})
	} else {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
	}

}

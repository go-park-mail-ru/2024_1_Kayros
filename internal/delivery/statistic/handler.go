package statistic

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"

	"github.com/mailru/easyjson"
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
	ucSession session.Usecase
	logger    *zap.Logger
	cfg       *config.Project
	metrics   *metrics.Metrics
}

func NewDeliveryLayer(ucQuizProps statistic.Usecase, ucUserProps user.Usecase, ucSessionProps session.Usecase, loggerProps *zap.Logger, cfgProps *config.Project, metrics   *metrics.Metrics) *Delivery {
	return &Delivery{
		ucQuiz:    ucQuizProps,
		ucUser:    ucUserProps,
		logger:    loggerProps,
		ucSession: ucSessionProps,
		cfg:       cfgProps,
		metrics: metrics,
	}
}

func (d *Delivery) GetStatistic(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email != "" {
		d.logger.Error(myerrors.CtxEmail.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	stats, err := d.ucQuiz.GetStatistic(r.Context())
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusUnauthorized)
		return
	}

	statisticArray := &dto.StatisticArray{Payload: dto.NewDtoStatistic(stats)} 
	functions.JsonResponse(w, statisticArray)
}

func (d *Delivery) GetQuestions(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)

	url := r.URL.Query().Get("url")
	qs := []*entity.Question{}
	var err error
	if url != "" {
		qs, err = d.ucQuiz.GetQuestionsOnFocus(r.Context(), url)
		if err != nil {
			d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
	}
	questionArray := &dto.QuestionArray{Payload: dto.QuestionReturn(qs)} 
	functions.JsonResponse(w, questionArray)
}

func (d *Delivery) AddAnswer(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)

	var qi []*dto.QuestionInput
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	qArray := dto.QuestionInputArray{Payload: qi}
	err = easyjson.Unmarshal(body, &qArray)
	if err != nil {
		d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
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
				w, err = functions.FlashCookie(r, w, d.ucSession, &d.cfg.Redis, d.metrics)
				if err != nil {
					d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
					functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
					return
				}
				functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
				return
			}
			functions.ErrorResponse(w, myerrors.InternalServerEn, http.StatusInternalServerError)
			return
		}
		for _, q := range qi {
			err = d.ucQuiz.Create(r.Context(), q.Id, q.Rating, strconv.Itoa(int(u.Id)))
			if err != nil {
				d.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				if errors.Is(err, myerrors.QuizAdd) {
					functions.ErrorResponse(w, myerrors.QuizAddRu, http.StatusInternalServerError)
				} else {
					functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
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
				functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
		}
		hasVoted = true
	}
	if hasVoted {
		functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Пользователь успешно проголосовал"})
	} else {
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
	}

}

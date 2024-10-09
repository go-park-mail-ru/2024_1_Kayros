package comment

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/usecase/comment"
	"2024_1_kayros/internal/utils/alias"

	"2024_1_kayros/internal/entity/dto"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Delivery struct {
	uc     comment.Usecase
	logger *zap.Logger
}

func NewDelivery(ucc comment.Usecase, loggerProps *zap.Logger) *Delivery {
	return &Delivery{
		uc:     ucc,
		logger: loggerProps,
	}
}

func (h *Delivery) CreateComment(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		return
	}

	var inputComment dto.InputComment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = easyjson.Unmarshal(body, &inputComment)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	com := entity.Comment{
		RestId: uint64(id),
		Text:   inputComment.Text,
		Rating: inputComment.Rating,
	}

	res, err := h.uc.CreateComment(r.Context(), com, email, inputComment.OrderId)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsCommentRelation) {
			functions.ErrorResponse(w, myerrors.FailCreateCommentRu, http.StatusInternalServerError)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	comDTO := dto.NewComment(res)
	functions.JsonResponse(w, comDTO)
}

func (h *Delivery) GetComments(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	comments, err := h.uc.GetCommentsByRest(r.Context(), alias.RestId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	if comments == nil {
		functions.ErrorResponse(w, myerrors.NoCommentsRu, http.StatusOK)
		return
	}

	commentArrayDTO := &dto.CommentArray{Payload: dto.NewCommentArray(comments)}
	functions.JsonResponse(w, commentArrayDTO)
}

func (h *Delivery) DeleteComment(w http.ResponseWriter, r *http.Request) {
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["com_id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		return
	}

	err = h.uc.DeleteComment(r.Context(), uint64(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsCommentRelation) {
			functions.ErrorResponse(w, myerrors.FailCreateCommentRu, http.StatusNotFound)
			return
		}
		functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	functions.JsonResponse(w, &dto.ResponseDetail{Detail: "Комментарий успешно удален"})
}

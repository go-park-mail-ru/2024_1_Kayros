package comment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

type InputId struct {
	Id uint64 `json:"id"`
}

type InputComment struct {
	Text   string `json:"text"`
	Rating uint32 `json:"rating"`
}

func (h *Delivery) CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		return
	}

	var inputComment *InputComment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &inputComment)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	com := entity.Comment{
		RestId: uint64(id),
		Text:   inputComment.Text,
		Rating: inputComment.Rating,
	}

	res, err := h.uc.CreateComment(r.Context(), com, email)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsCommentRelation) {
			w = functions.ErrorResponse(w, myerrors.FailCreateCommentRu, http.StatusInternalServerError)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	comDTO := dto.NewComment(res)
	w = functions.JsonResponse(w, comDTO)
}

func (h *Delivery) GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	comments, err := h.uc.GetCommentsByRest(r.Context(), alias.RestId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	if comments == nil {
		w = functions.ErrorResponse(w, myerrors.NoCommentsRu, http.StatusOK)
		return
	}

	commentArrayDTO := dto.NewCommentArray(comments)
	w = functions.JsonResponse(w, commentArrayDTO)
}

func (h *Delivery) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["com_id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.NotFoundRu, http.StatusNotFound)
		return
	}

	err = h.uc.DeleteComment(r.Context(), uint64(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.SqlNoRowsCommentRelation) {
			w = functions.ErrorResponse(w, myerrors.FailCreateCommentRu, http.StatusNotFound)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
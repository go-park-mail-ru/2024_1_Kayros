package comment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"2024_1_kayros/internal/usecase/comment"
	"2024_1_kayros/internal/utils/alias"

	"2024_1_kayros/internal/entity/dto"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type CommentHandler struct {
	uc     comment.Usecase
	logger *zap.Logger
}

func NewCommentHandler(ucc comment.Usecase, loggerProps *zap.Logger) *CommentHandler {
	return &CommentHandler{
		uc:     ucc,
		logger: loggerProps,
	}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)
	email := functions.GetCtxEmail(r)
	if email == "" {
		h.logger.Error(myerrors.UnauthorizedEn.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
		return
	}

	var comDTO *dto.Comment

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
	err = json.Unmarshal(body, &comDTO)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	com := dto.NewCommentFromDTO(comDTO)

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
	comDTO = dto.NewComment(res)
	w = functions.JsonResponse(w, comDTO)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestId := functions.GetCtxRequestId(r)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["rest_id"])
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		w = functions.ErrorResponse(w, myerrors.BadCredentialsRu, http.StatusBadRequest)
		return
	}

	comments, err := h.uc.GetCommentsByRest(r.Context(), alias.RestId(id))
	if err != nil {
		h.logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		if errors.Is(err, myerrors.NoComments) {
			w = functions.ErrorResponse(w, myerrors.NoCommentsRu, http.StatusOK)
			return
		}
		w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
		return
	}

	commentArrayDTO := dto.NewCommentArray(comments)
	w = functions.JsonResponse(w, commentArrayDTO)
}

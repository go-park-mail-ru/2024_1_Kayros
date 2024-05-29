package functions

import (
	"2024_1_kayros/internal/entity/dto"
	"net/http"

	"github.com/mailru/easyjson"
)

func ErrorResponse(w http.ResponseWriter, responseError error, codeStatus int) {
	w.Header().Add("Content-Type", "application/json")
	errObject := &dto.ResponseDetail{Detail: responseError.Error()}
	body, err := easyjson.Marshal(errObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(codeStatus)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

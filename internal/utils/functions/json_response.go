package functions

import (
	"encoding/json"
	"net/http"

	"2024_1_kayros/internal/utils/myerrors"
)

func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	body, err := json.Marshal(data)
	if err != nil {
		w = ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	if err != nil {
		w = ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}
	w.WriteHeader(http.StatusOK)
	return w
}

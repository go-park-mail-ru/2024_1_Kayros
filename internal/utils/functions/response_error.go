package functions

import (
	"encoding/json"
	"net/http"

	"2024_1_kayros/internal/utils/myerrors"
)

// ErrorResponse формирует ответ с ошибкой от хендлеров
func ErrorResponse(w http.ResponseWriter, messageError string, codeStatus int) http.ResponseWriter {
	errObject := map[string]string{"detail": messageError}
	responseBody, err := json.Marshal(errObject)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(myerrors.UnexpectedServerError))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return w
	}
	w.WriteHeader(codeStatus)
	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return w
}
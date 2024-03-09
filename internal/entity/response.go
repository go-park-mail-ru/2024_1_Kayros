package entity

import (
	"encoding/json"
	"net/http"
)

type ErrorObject struct {
	Detail string `json:"detail"`
}

// ErrorResponse формирует ответ с ошибкой от хендлеров
func ErrorResponse(w http.ResponseWriter, message string, code int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	errObject := ErrorObject{Detail: message}
	responseBody, errSerialization := json.Marshal(errObject)
	if errSerialization != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMessageBody, _ := json.Marshal("Ошибка формирования ответа")
		_, errWriteErrorResponseBody := w.Write(errMessageBody)
		if errWriteErrorResponseBody != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return w
	}
	w.WriteHeader(code)
	_, errWriteResponseBody := w.Write(responseBody)
	if errWriteResponseBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return w
}
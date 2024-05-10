package functions

import (
	"encoding/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, responseError error, codeStatus int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	errObject := map[string]string{"detail": responseError.Error()}
	body, err := json.Marshal(errObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}

	w.WriteHeader(codeStatus)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return w
}
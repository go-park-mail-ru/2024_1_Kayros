package functions

import (
	"encoding/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, messageError string, codeStatus int) http.ResponseWriter {
	errObject := map[string]string{"detail": messageError}
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

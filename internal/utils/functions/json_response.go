package functions

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}
	return w
}

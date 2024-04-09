package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}
	return w
}

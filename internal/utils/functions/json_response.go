package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
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

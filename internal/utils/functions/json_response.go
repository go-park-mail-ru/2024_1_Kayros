package functions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"2024_1_kayros/internal/utils/myerrors"
)

func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	body, err := json.Marshal(data)
	fmt.Println(err)
	if err != nil {
		w = ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(body)
	fmt.Println(err)
	if err != nil {
		w = ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}
	return w
}

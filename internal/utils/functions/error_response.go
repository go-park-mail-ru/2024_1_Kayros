package functions

import (
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, messageError string, codeStatus int) http.ResponseWriter {
	errObject := map[string]string{"detail": messageError}
	w = JsonResponse(w, errObject)
	//w.WriteHeader(codeStatus)
	return w
}

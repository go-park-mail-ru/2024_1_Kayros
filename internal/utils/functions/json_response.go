package functions

import (
	"net/http"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

func JsonResponse(w http.ResponseWriter, data interface{}) {
	logger := zap.Logger{}
	var err error
	_, _, err = easyjson.MarshalToHTTPResponseWriter(data.(easyjson.Marshaler), w)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

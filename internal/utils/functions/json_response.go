package functions

import (
	"2024_1_kayros/internal/entity/dto"
	"net/http"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)


func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	logger := zap.Logger{}
	var err error
	switch dtoData := data.(type) {
	case dto.Comment:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Food:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.FoodInOrder:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Order:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.ShortOrder:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.FullAddress:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Promocode:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.QuestionInput:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Question:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Category:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.RestaurantAndFood:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Restaurant:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Statistic:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.UserUpdate:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.UserSignUp:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.UserSignIn:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.UserGet:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Address:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	case dto.Passwords:
		_, _, err = easyjson.MarshalToHTTPResponseWriter(dtoData, w)
	}
	if err != nil {
		logger.Error(err.Error())
		return w
	}

	return w
}

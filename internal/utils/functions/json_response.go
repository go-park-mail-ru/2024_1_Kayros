package functions

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)


func JsonResponse(w http.ResponseWriter, data interface{}) http.ResponseWriter {
	logger := zap.Logger{}
	var err error
	
	switch reflect.TypeOf(data).String() {
		case "*dto.Address", "*dto.ResponseDetail", "*dto.UserGet",
		"*dto.StatisticArray", "*dto.QuestionArray", "*dto.RestaurantAndFoodArray",
		"*dto.RestaurantArray", "*dto.RestaurantAndFood", "*dto.CategoryArray",
		"*dto.ResponseUrlPay", "*dto.Order", "*dto.ShortOrderArray", 
		"*dto.PayedOrderInfo", "*dto.Promo", "*dto.Comment", 
		"*dto.CommentArray":
			_, _, err = easyjson.MarshalToHTTPResponseWriter(data.(easyjson.Marshaler), w)
	default:
		fmt.Println("NO MATCH")
	}
	if err != nil {
		logger.Error(err.Error())
		return w
	}

	return w
}

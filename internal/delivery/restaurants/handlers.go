package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"2024_1_kayros/internal/entity"
)

var rests = []entity.Restaurant{
	{1, "Пицца 22 см", "Пиццерия с настоящей неаполитанской пиццей", "1.jpg"},
	{2, "Bro&N", "Ресторан классической итальянской кухни", "2.jpg"},
	{3, "#FARШ", "Сеть бургерных с сочным мясом от \"Мираторга\"", "3.jpg"},
	{4, "Loona", "Итальянскую классику в современном прочтении", "4.jpg"},
	{5, "Pino", "Обширное интернациональное меню", "5.jpg"},
	{6, "Sage", "Авторская евпропейская кухня с акцентом на мясо и рыбу", "6.jpg"},
	{7, "TECHNIKUM", "Современное гастробистро с нескучной едой", "7.jpg"},
}

func RestaurantList(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(rests)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))

	_, err = w.Write(body)
	if err != nil {
		fmt.Printf("Write failed: %v", err)
	}
}

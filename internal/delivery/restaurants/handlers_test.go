package delivery

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Response   string
	StatusCode int
}

func TestRests(t *testing.T) {
	casee := TestCase{
		Response:   `[{"id":1,"name":"Пицца 22 см","description":"Пиццерия с настоящей неаполитанской пиццей","img_url":"assets/mocks/restaurants/1.jpg"},{"id":2,"name":"Bro\u0026N","description":"Ресторан классической итальянской кухни","img_url":"assets/mocks/restaurants/2.jpg"},{"id":3,"name":"#FARШ","description":"Сеть бургерных с сочным мясом от \"Мираторга\"","img_url":"assets/mocks/restaurants/3.jpg"},{"id":4,"name":"Loona","description":"Итальянскую классику в современном прочтении","img_url":"assets/mocks/restaurants/4.jpg"},{"id":5,"name":"Pino","description":"Обширное интернациональное меню","img_url":"assets/mocks/restaurants/5.jpg"},{"id":6,"name":"Sage","description":"Авторская евпропейская кухня с акцентом на мясо и рыбу","img_url":"assets/mocks/restaurants/6.jpg"},{"id":7,"name":"TECHNIKUM","description":"Современное гастробистро с нескучной едой","img_url":"assets/mocks/restaurants/7.jpg"}]`,
		StatusCode: http.StatusOK,
	}
	rests := NewRestaurantStore()
	url := "https://resto-go.ru/restaurants"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	rests.RestaurantList(w, req)
	if w.Code != casee.StatusCode {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, casee.StatusCode)
	}
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if bodyStr != casee.Response {
		t.Errorf("wrong Response: got %+v, expected %+v", bodyStr, casee.Response)
	}
}

func TestError(t *testing.T) {
	casee := TestCase{
		Response:   ``,
		StatusCode: http.StatusOK,
	}
	fmt.Printf("", casee)
	rests := NewRestaurantStore()
	url := "https://resto-go.ru/restaurants"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	rests.RestaurantList(w, req)
	if w.Code != casee.StatusCode {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, casee.StatusCode)
	}
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if bodyStr != casee.Response {
		t.Errorf("wrong Response: got %+v, expected %+v", bodyStr, casee.Response)
	}
}

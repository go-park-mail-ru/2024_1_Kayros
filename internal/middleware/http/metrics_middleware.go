package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/utils/recorder"
)

func Metrics(handler http.Handler, m *metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start metrics work
		start := time.Now()
		// init custom http.ResponseWriter
		rec := recorder.NewResponseWriter(w)
		// call handler
		handler.ServeHTTP(rec, r)
		// collect metrics data
		method := r.Method
		url := r.URL.String()

		status := strconv.Itoa(rec.StatusCode)
		classificatedURL := classificateUrlForMetrics(url)
		// fill metrics data
		m.TotalNumberOfRequests.Inc()
		m.NumberOfSpecificRequests.WithLabelValues(method, classificatedURL, status).Inc()
		m.RequestTime.WithLabelValues(method, classificatedURL, status).Observe(float64(time.Since(start).Milliseconds()))
	})
}

func classificateUrlForMetrics(url string) string {
	urlSplit := strings.Split(url, "/")
	var numIndexInURL []int
	for index, part := range urlSplit {
		_, err := strconv.Atoi(part)
		if err == nil {
			numIndexInURL = append(numIndexInURL, index)
		}
	}

	needToClassificate := false
	indexChangle := 0
	if len(numIndexInURL) != 0 {
		needToClassificate = true
		indexChangle = numIndexInURL[len(numIndexInURL)-1]
	}
	newUrl := ""
	for index, part := range urlSplit {
		if index != 0 {
			newUrl += "/"
		}
		if index == indexChangle && needToClassificate {
			newUrl += "id"
		} else {
			newUrl += part
		}
	}
	return newUrl
}

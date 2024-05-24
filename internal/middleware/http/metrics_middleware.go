package http

import (
	"net/http"
	"strconv"
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
		// fill metrics data
		m.TotalNumberOfRequests.Inc()
		m.NumberOfSpecificRequests.WithLabelValues(method, url, status).Inc()
		m.RequestTime.WithLabelValues(method, url, status).Observe(float64(time.Since(start).Milliseconds()))
	})
}
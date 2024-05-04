package middleware

import (
	"net/http"
	"time"

	metrics "2024_1_kayros"
)

func Metrics(handler http.Handler, m *metrics.Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		url := r.URL.String()
		handler.ServeHTTP(w, r)
		m.Hits.WithLabelValues(url).Inc()
		m.Duration.WithLabelValues(url, method).Observe(time.Since(start).Seconds())
	})
}

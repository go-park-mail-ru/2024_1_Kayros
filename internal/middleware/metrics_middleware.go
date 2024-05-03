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
		m.Hits.WithLabelValues(method, r.URL.String()).Inc()
		m.Duration.WithLabelValues(method, r.URL.String()).Observe(time.Since(start).Seconds())
		handler.ServeHTTP(w, r)
	})
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)
)

func main() {
	prometheus.MustRegister(requestsTotal, requestDuration)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		method := r.Method
		elapsed := time.Since(start).Seconds()
		requestsTotal.WithLabelValues(method).Inc()
		requestDuration.WithLabelValues(method).Observe(elapsed)
	}()

	// Ваш код обработки запроса здесь
	fmt.Fprintf(w, "Hello, World!")
}

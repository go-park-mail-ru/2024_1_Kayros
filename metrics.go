package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)
)

func main() {
	prometheus.MustRegister(hits, requestDuration)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			method := r.Method
			elapsed := time.Since(start).Seconds()
			hits.WithLabelValues(method, r.URL.String()).Inc()
			requestDuration.WithLabelValues(method).Observe(elapsed)
		}()
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("metrics server can't be started")
	}
}

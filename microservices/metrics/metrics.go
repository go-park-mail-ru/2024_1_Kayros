package microservice_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MicroserviceMetrics struct {
	TotalNumberOfRequests prometheus.Counter
	RequestTime           *prometheus.HistogramVec
	DatabaseDuration      *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer, namespace string) *MicroserviceMetrics {
	m := &MicroserviceMetrics{
		// business metrics
		// total number of hits
		TotalNumberOfRequests: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "number_of_requests",
				Help:      "number of requests",
			},
		),
		// request time can be filtered by request method, url path and response status
		RequestTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "time_of_request",
				Help:      "HTTP request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"status"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_duration_ms",
				Help:      "Database request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),
	}
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(m.TotalNumberOfRequests)
	reg.MustRegister(m.RequestTime)
	reg.MustRegister(m.DatabaseDuration)
	return m
}

const (
	SELECT = "SELECT"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	INSERT = "INSERT"
	REDIS  = "REDIS"
)

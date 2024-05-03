package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Hits     *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Hits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hits",
				Help: "Number of hits.",
			},
			[]string{"path", "status"},
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "duration",
				Help: "Duration of request",
			},
			[]string{"path", "status"},
		),
	}
	reg.MustRegister(m.Hits, m.Duration)
	return m
}

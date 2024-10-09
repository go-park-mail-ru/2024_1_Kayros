package repo

import (
	metrics "2024_1_kayros/microservices/metrics"
	"database/sql"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testFixtures struct {
	repo Rest
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func setUp(t *testing.T) testFixtures {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	m := &metrics.MicroserviceMetrics{
		// business metrics
		// total number of hits
		TotalNumberOfRequests: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "",
				Name:      "number_of_requests",
				Help:      "number of requests",
			},
		),
		// request time can be filtered by request method, url path and response status
		RequestTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "",
				Name:      "time_of_request",
				Help:      "HTTP request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"status"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "",
				Name:      "database_duration_ms",
				Help:      "Database request duration in milliseconds",
				Buckets:   []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),
	}
	repo := NewRestLayer(db, m)
	return testFixtures{
		repo: repo,
		db:   db,
		mock: mock,
	}
}

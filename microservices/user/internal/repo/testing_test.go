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
	repo Repo
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func setUp(t *testing.T, namespace string) testFixtures {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := NewLayer(db, &metrics.MicroserviceMetrics{
		RequestTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:    "time_of_request",
				Help:    "HTTP request duration in milliseconds",
				Buckets: []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"status"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:    "database_duration_ms",
				Help:    "Database request duration in milliseconds",
				Buckets: []float64{10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"operation"},
		),
	})
	return testFixtures{
		repo: repo,
		db:   db,
		mock: mock,
	}
}
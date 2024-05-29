package repo

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testFixtures struct {
	repo Repo
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func setUp(t *testing.T) testFixtures {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := NewLayer(db)
	return testFixtures{
		repo: repo,
		db:   db,
		mock: mock,
	}
}

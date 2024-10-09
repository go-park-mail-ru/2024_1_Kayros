package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"2024_1_kayros/gen/go/rest"
	"2024_1_kayros/internal/utils/myerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetAll(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t)
	defer s.db.Close()

	t.Run("db error", func(t *testing.T) {

		s.mock.
			ExpectQuery("SELECT id, name, short_description, address, img_url FROM restaurant ORDER BY rating DESC").
			WillReturnError(fmt.Errorf("db_error"))

		_, err := s.repo.GetAll(ctx)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "name", "short_description", "address", "img_url"})
		restExpected := &rest.RestList{
			Rest: []*rest.Rest{
				{
					Id:               1,
					Name:             "a",
					ShortDescription: "a",
					Address:          "a",
					ImgUrl:           "a",
				},
			},
		}
		for _, r := range restExpected.Rest {
			rows = rows.AddRow(r.Id, r.Name, r.ShortDescription, r.Address, r.ImgUrl)
		}

		s.mock.
			ExpectQuery("SELECT id, name, short_description, address, img_url FROM restaurant ORDER BY rating DESC").
			WillReturnRows(rows)
		rests, err := s.repo.GetAll(ctx)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, restExpected, rests)

	})
}

func TestGetById(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t)
	defer s.db.Close()

	t.Run("db error", func(t *testing.T) {

		rId := &rest.RestId{
			Id: 1,
		}

		s.mock.
			ExpectQuery("SELECT id, name, long_description, address, img_url, rating, comment_count FROM restaurant WHERE").
			WithArgs(rId.Id).
			WillReturnError(fmt.Errorf("db_error"))

		_, err := s.repo.GetById(ctx, rId)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("no rows", func(t *testing.T) {

		rId := &rest.RestId{
			Id: 1,
		}

		s.mock.
			ExpectQuery("SELECT id, name, long_description, address, img_url, rating, comment_count FROM restaurant WHERE").
			WithArgs(rId.Id).
			WillReturnError(sql.ErrNoRows)

		_, err := s.repo.GetById(ctx, rId)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, myerrors.SqlNoRowsRestaurantRelation, err)
	})

}

func TestGetByFilter(t *testing.T) {
	s := setUp(t)
	t.Parallel()
	ctx := context.Background()

	t.Run("db error", func(t *testing.T) {
		id := &rest.Id{Id: 1}

		s.mock.
			ExpectQuery("SELECT r.id, r.name, r.short_description, r.img_url FROM restaurant as r \n\t\t\t\tJOIN rest_categories AS rc ON r.id=rc.restaurant_id WHERE").
			WithArgs(id.Id).
			WillReturnError(fmt.Errorf("db_error"))

		_, err := s.repo.GetByFilter(ctx, id)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		id := &rest.Id{Id: 1}

		rows := sqlmock.NewRows([]string{"id", "name", "short_description", "img_url"})
		restExpected := &rest.RestList{
			Rest: []*rest.Rest{
				{
					Id:               1,
					Name:             "a",
					ShortDescription: "a",
					ImgUrl:           "a",
				},
			},
		}
		for _, r := range restExpected.Rest {
			rows = rows.AddRow(r.Id, r.Name, r.ShortDescription, r.ImgUrl)
		}

		s.mock.
			ExpectQuery("SELECT r.id, r.name, r.short_description, r.img_url FROM restaurant as r \n\t\t\t\tJOIN rest_categories AS rc ON r.id=rc.restaurant_id WHERE").
			WithArgs(id.Id).
			WillReturnRows(rows)
		rests, err := s.repo.GetByFilter(ctx, id)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, restExpected, rests)

	})

}

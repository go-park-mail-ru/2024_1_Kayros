package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t, "TestGetByEmail")
	defer s.db.Close()

	t.Run("internal db error", func(t *testing.T) {
		email := &user.Email{Email: "aaa@aa.aa"}

		s.mock.ExpectQuery(`SELECT id, name, email, COALESCE\(phone, ''\), password, COALESCE\(address, ''\), img_url, COALESCE\(card_number, ''\) FROM "user" WHERE email = \$1`).
			WithArgs(email.Email).
			WillReturnError(fmt.Errorf("db_error"))

		_, err := s.repo.GetByEmail(ctx, email)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		email := &user.Email{Email: "aaa@aa.aa"}

		s.mock.ExpectQuery(`SELECT id, name, email, COALESCE\(phone, ''\), password, COALESCE\(address, ''\), img_url, COALESCE\(card_number, ''\) FROM "user" WHERE email = \$1`).
			WithArgs(email.Email).
			WillReturnError(sql.ErrNoRows)

		res, err := s.repo.GetByEmail(ctx, email)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, myerrors.SqlNoRowsUserRelation, err)
		assert.Nil(t, res)
	})

	t.Run("ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "address", "img_url", "card_number"})
		email :=  "aaa@aa.aa"
		u := &entity.User{
			Id:         1,
			Name:       "aaa",
			Email:      "aaa@aa.aa",
			Phone:      "88005553535",
			Password:   "qwerty123",
			Address:    "aaa",
			ImgUrl:     "aaaa.a",
			CardNumber: "1234",
		}
		rows = rows.AddRow(u.Id, u.Name, u.Email, u.Phone, u.Password, u.Address, u.ImgUrl, u.CardNumber)

		s.mock.ExpectQuery(`SELECT id, name, email, COALESCE\(phone, ''\), password, COALESCE\(address, ''\), img_url, COALESCE\(card_number, ''\) FROM "user" WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(rows)
		uReturn, err := s.repo.GetByEmail(ctx, &user.Email{Email: email})

		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Equal(t, entity.ConvertEntityUserIntoProtoUser(u), uReturn)
	})
}

func TestDeleteByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t, "TestDeleteByEmail")
	defer s.db.Close()

	t.Run("db internal error", func (t *testing.T)  {
		email := &user.Email{Email: "aaa@aaa.aa"}
		s.mock.ExpectExec(`DELETE FROM "user" WHERE email = \$1`).WithArgs(email.GetEmail()).WillReturnError(fmt.Errorf("db_error"))
		
		err := s.repo.DeleteByEmail(ctx, email)
		assert.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})

	t.Run("db error no rows affected", func (t *testing.T)  {
		email := &user.Email{Email: "aaa@aaa.aa"}
		s.mock.ExpectExec(`DELETE FROM "user" WHERE email = \$1`).WithArgs(email.GetEmail()).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("rows affected error")))
		
		err := s.repo.DeleteByEmail(ctx, email)
		assert.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, err, myerrors.SqlNoRowsUserRelationAffected)
	})

	t.Run("db no rows affected", func (t *testing.T)  {
		email := &user.Email{Email: "aaa@aaa.aa"}
		s.mock.ExpectExec(`DELETE FROM "user" WHERE email = \$1`).WithArgs(email.GetEmail()).WillReturnResult(sqlmock.NewResult(0, 0))
		
		err := s.repo.DeleteByEmail(ctx, email)
		assert.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, err, myerrors.SqlNoRowsUserRelation)
	})

	t.Run("ok", func (t *testing.T)  {
		email := &user.Email{Email: "aaa@aaa.aa"}
		s.mock.ExpectExec(`DELETE FROM "user" WHERE email = \$1`).WithArgs(email.GetEmail()).WillReturnResult(sqlmock.NewResult(3, 1))
		
		err := s.repo.DeleteByEmail(ctx, email)
		assert.NoError(t, s.mock.ExpectationsWereMet())
		assert.Nil(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()
	s := setUp(t, "TestCreate")
	defer s.db.Close()
	ctx := context.Background()

	t.Run("db internal error", func(t *testing.T) {
		u := &user.User{
			Name:     "a",
			Email:    "aaa@aa.aa",
			Password: "a",
		}

		s.mock.ExpectExec(`INSERT INTO "user" \(name, email, password, address, img_url, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`).
			WithArgs(u.GetName(), u.GetEmail(), u.GetPassword(), functions.MaybeNullString(u.GetAddress()), functions.MaybeNullString(u.GetImgUrl()),  sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("db_error"))

		err := s.repo.Create(ctx, u)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)

	})

	t.Run("db error no rows affected", func(t *testing.T) {
		u := &user.User{
			Name:     "a",
			Email:    "aaa@aa.aa",
			Password: "a",
		}

		s.mock.ExpectExec(`INSERT INTO "user" \(name, email, password, address, img_url, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`).
			WithArgs(u.GetName(), u.GetEmail(), u.GetPassword(), functions.MaybeNullString(u.GetAddress()), functions.MaybeNullString(u.GetImgUrl()),  sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("rows affected error")))

		err := s.repo.Create(ctx, u)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, myerrors.SqlNoRowsUserRelationAffected, err)
	})

	t.Run("no rows affected", func(t *testing.T) {
		u := &user.User{
			Name:     "a",
			Email:    "aaa@aa.aa",
			Password: "a",
		}

		s.mock.ExpectExec(`INSERT INTO "user" \(name, email, password, address, img_url, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`).
			WithArgs(u.GetName(), u.GetEmail(), u.GetPassword(), functions.MaybeNullString(u.GetAddress()), functions.MaybeNullString(u.GetImgUrl()),  sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repo.Create(ctx, u)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, myerrors.SqlNoRowsUserRelation, err)

	})

	t.Run("ok", func(t *testing.T) {
		u := &user.User{
			Name:     "a",
			Email:    "aaa@aa.aa",
			Password: "a",
		}

		s.mock.ExpectExec(`INSERT INTO "user" \(name, email, password, address, img_url, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`).
			WithArgs(u.GetName(), u.GetEmail(), u.GetPassword(), functions.MaybeNullString(u.GetAddress()), functions.MaybeNullString(u.GetImgUrl()),  sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.repo.Create(ctx, u)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.NoError(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t, "TestUpdate")

	t.Run("db internal error", func(t *testing.T) {
		email := &user.Email{Email: "aaa@aa.aa"}

		user := &user.UpdateUserData{
			UpdateInfo: &user.User{
				Name:     "a",
				Email:    "aaa@aa.aa",
				Password: "a",
			},
			Email: "aaa@aaa.aa",
			FileData: []byte{},
			FileName: "123.svg",
			FileSize: 123,
		}

		userData := user.UpdateInfo
		s.mock.ExpectExec(`UPDATE "user" SET name = \$1, email = \$2, phone = \$3, img_url = \$4, password = \$5, card_number = \$6, 
		address = \$7, updated_at = \$8 WHERE email = \$9`).
		WithArgs(userData.GetName(), userData.GetEmail(), functions.MaybeNullString(userData.GetPhone()), 
		functions.MaybeNullString(userData.GetImgUrl()), userData.GetPassword(), functions.MaybeNullString(userData.GetCardNumber()), functions.MaybeNullString(userData.GetAddress()), 
		"2006-01-02 15:04:05-07:00", email.GetEmail()).WillReturnError(fmt.Errorf("db error"))

		err := s.repo.Update(ctx, user)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)
	})
}

func TestGetAddressByUnauthId(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := setUp(t, "TestGetAddressByUnauthId")
	defer s.db.Close()

	t.Run("db internal error", func(t *testing.T) {
		id := &user.UnauthId{UnauthId: "2ede4ef8-0bd0-4462-9350-220c44729072"}

		s.mock.ExpectQuery(`SELECT address FROM unauth_address WHERE unauth_id = \$1`).
			WithArgs(id.UnauthId).
			WillReturnError(fmt.Errorf("db_error"))

		_, err := s.repo.GetAddressByUnauthId(ctx, id)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Error(t, err)

	})

	t.Run("not found", func(t *testing.T) {
		id := &user.UnauthId{UnauthId: "a"}

		s.mock.ExpectQuery(`SELECT address FROM unauth_address WHERE unauth_id = \$1`).
			WithArgs(id.UnauthId).
			WillReturnError(sql.ErrNoRows)

		_, err := s.repo.GetAddressByUnauthId(ctx, id)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.Equal(t, myerrors.SqlNoRowsUnauthAddressRelation, err)
	})


	t.Run("null address", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"address"}).AddRow(nil)
		id := &user.UnauthId{UnauthId: "a"}

		s.mock.ExpectQuery(`SELECT address FROM unauth_address WHERE unauth_id = \$1`).
			WithArgs(id.UnauthId).
			WillReturnRows(rows)
		address, err := s.repo.GetAddressByUnauthId(ctx, id)
		
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Nil(t, address)
	})

	t.Run("ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"address"})
		rows = rows.AddRow("a")
		id := &user.UnauthId{UnauthId: "a"}

		s.mock.ExpectQuery(`SELECT address FROM unauth_address WHERE unauth_id = \$1`).
			WithArgs(id.UnauthId).
			WillReturnRows(rows)
		address, err := s.repo.GetAddressByUnauthId(ctx, id)
		require.NoError(t, s.mock.ExpectationsWereMet())
		assert.NoError(t, err)
		assert.Equal(t, &user.Address{Address: "a"}, address)
	})
}

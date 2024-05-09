package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/internal/entity"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/gen/go/user"
)

type Repo interface {
	GetByEmail(ctx context.Context, email *user.Email) (*user.User, error)
	DeleteByEmail(ctx context.Context, email *user.Email) error
	Create(ctx context.Context, u *user.User) error
	Update(ctx context.Context, data *user.UpdateUserData) error
	GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error)
	UpdateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error
	CreateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error
}

type Layer struct {
	database *sql.DB
}

func NewLayer(db *sql.DB) Repo {
	return &Layer{
		database: db,
	}
}

func (repo *Layer) GetByEmail(ctx context.Context, email *user.Email) (*user.User, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, '')  FROM "user" WHERE email = $1`, email.GetEmail())
	u := entity.User{}
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Address, &u.ImgUrl, &u.CardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return &user.User{}, err
	}
	return entity.ConvertEntityUserIntoProtoUser(&u), nil
}


func (repo *Layer) DeleteByEmail(ctx context.Context, email *user.Email) error {
	row, err := repo.database.ExecContext(ctx, `DELETE FROM "user" WHERE email = $1`, email.GetEmail())
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}

func (repo *Layer) Create(ctx context.Context, u *user.User) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	row, err := repo.database.ExecContext(ctx,
		`INSERT INTO "user" (name, email, password, address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		u.GetName(), u.GetEmail(), u.GetPassword(), functions.MaybeNullString(u.GetAddress()), timeNow, timeNow)
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}

func (repo *Layer) Update(ctx context.Context, data *user.UpdateUserData) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	userData := data.GetUpdateInfo()
	row, err := repo.database.ExecContext(ctx,
		`UPDATE "user" SET name = $1, email = $2, phone = $3, img_url = $4, password = $5, card_number = $6, 
                  address = $7, updated_at = $8 WHERE email = $9`,
		userData.GetName(), userData.GetEmail(), functions.MaybeNullString(userData.GetPhone()), userData.GetImgUrl(),
		userData.GetPassword(), functions.MaybeNullString(userData.GetCardNumber()), functions.MaybeNullString(userData.GetAddress()), timeNow, data.GetEmail())
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUserRelation
	}
	return nil
}

func (repo *Layer) GetAddressByUnauthId(ctx context.Context, id *user.UnauthId) (*user.Address, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT address FROM unauth_address WHERE unauth_id = $1`, id.GetUnauthId())
	var address sql.NullString
	err := row.Scan(&address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user.Address{}, myerrors.SqlNoRowsUnauthAddressRelation
		}
		return &user.Address{}, err
	}
	if !address.Valid {
		return &user.Address{}, nil
	}
	return &user.Address{Address: address.String}, nil
}

func (repo *Layer) UpdateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error {
	row, err := repo.database.ExecContext(ctx, `UPDATE unauth_address SET address = $1 WHERE unauth_id= $2`, functions.MaybeNullString(data.GetAddress()), data.GetUnauthId())
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUnauthAddressRelation
	}
	return nil
}

func (repo *Layer) CreateAddressByUnauthId(ctx context.Context, data *user.AddressDataUnauth) error {
	row, err := repo.database.ExecContext(ctx, `INSERT INTO unauth_address (unauth_id, address) VALUES ($1, $2)`, data.GetUnauthId(), data.GetAddress())
	if err != nil {
		return err
	}
	numRows, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		return myerrors.SqlNoRowsUnauthAddressRelation
	}
	return nil
}
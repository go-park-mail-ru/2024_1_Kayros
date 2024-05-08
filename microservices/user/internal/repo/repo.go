package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"2024_1_kayros/internal/entity"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	userv1 "2024_1_kayros/microservices/user/proto"
)

type Repo interface {
	GetByEmail(ctx context.Context, email *userv1.Email) (*userv1.User, error)
	DeleteByEmail(ctx context.Context, email *userv1.Email) error
	Create(ctx context.Context, u *userv1.User) error
	Update(ctx context.Context, data *userv1.UpdateUserData) error
	GetAddressByUnauthId(ctx context.Context, id *userv1.UnauthId) (*userv1.Address, error)
	UpdateAddressByUnauthId(ctx context.Context, data *userv1.AddressDataUnauth) error
	CreateAddressByUnauthId(ctx context.Context, data *userv1.AddressDataUnauth) error
}

type Layer struct {
	database *sql.DB
}

func NewLayer(db *sql.DB) Repo {
	return &Layer{
		database: db,
	}
}

func (repo *Layer) GetByEmail(ctx context.Context, email *userv1.Email) (*userv1.User, error) {
	fmt.Printf("%v", email)
	row := repo.database.QueryRowContext(ctx,
		`SELECT id, name, email, COALESCE(phone, ''), password, COALESCE(address, ''), img_url, COALESCE(card_number, '')  FROM "user" WHERE email = $1`, email.GetEmail())
	user := entity.User{}
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Address, &user.ImgUrl, &user.CardNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return nil, err
	}
	return cnvUserIntoUserV1(&user), nil
}

func cnvUserIntoUserV1 (u *entity.User) *userv1.User {
	return &userv1.User{
		Id: u.Id,
		Name: u.Name,
		Phone: u.Phone,
		Email: &userv1.Email{Email: u.Email},
		Address: &userv1.Address{Address: u.Address},
		ImgUrl: u.ImgUrl,
		CardNumber: u.CardNumber,
		Password: &userv1.Password{Password: u.Password},
	}
}

func (repo *Layer) DeleteByEmail(ctx context.Context, email *userv1.Email) error {
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

func (repo *Layer) Create(ctx context.Context, u *userv1.User) error {
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	fmt.Printf("%v", u)
	row, err := repo.database.ExecContext(ctx,
		`INSERT INTO "user" (name, email, password, address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		u.GetName(), u.GetEmail().GetEmail(), u.GetPassword().GetPassword(), functions.MaybeNullString(u.GetAddress().GetAddress()), timeNow, timeNow)
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

func (repo *Layer) Update(ctx context.Context, data *userv1.UpdateUserData) error {
	fmt.Printf("%v", data)
	timeNow := time.Now().UTC().Format(cnst.Timestamptz)
	userData := data.GetUpdateInfo()
	row, err := repo.database.ExecContext(ctx,
		`UPDATE "user" SET name = $1, email = $2, phone = $3, img_url = $4, password = $5, card_number = $6, 
                  address = $7, updated_at = $8 WHERE email = $9`,
		userData.GetName(), userData.GetEmail().GetEmail(), functions.MaybeNullString(userData.GetPhone()), userData.GetImgUrl(),
		userData.GetPassword().GetPassword(), functions.MaybeNullString(userData.GetCardNumber()), functions.MaybeNullString(userData.GetAddress().GetAddress()), timeNow, data.GetEmail().GetEmail())
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

func (repo *Layer) GetAddressByUnauthId(ctx context.Context, id *userv1.UnauthId) (*userv1.Address, error) {
	row := repo.database.QueryRowContext(ctx,
		`SELECT address  FROM unauth_address WHERE unauth_id = $1`, id.GetUnauthId())
	var address sql.NullString
	err := row.Scan(&address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsUnauthAddressRelation
		}
		return nil, err
	}
	if !address.Valid {
		return nil, nil
	}
	return &userv1.Address{Address: address.String}, nil
}

func (repo *Layer) UpdateAddressByUnauthId(ctx context.Context, data *userv1.AddressDataUnauth) error {
	row, err := repo.database.ExecContext(ctx, `UPDATE unauth_address SET address = $1 WHERE unauth_id= $2`, functions.MaybeNullString(data.GetAddress().GetAddress()), data.GetId().GetUnauthId())
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

func (repo *Layer) CreateAddressByUnauthId(ctx context.Context, data *userv1.AddressDataUnauth) error {
	row, err := repo.database.ExecContext(ctx, `INSERT INTO unauth_address (unauth_id, address) VALUES ($1, $2)`, data.GetId().GetUnauthId(), data.GetAddress().GetAddress())
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
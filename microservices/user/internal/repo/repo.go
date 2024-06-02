package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/entity"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	metrics "2024_1_kayros/microservices/metrics"
)

//go:generate mockgen -source ./repo.go -destination=./mocks/repo.go -package=mock_repo
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
	metrics  *metrics.MicroserviceMetrics
	stmt     map[string]*sql.Stmt  // key - name of sql method, value - prepared statement 
}

func NewLayer(db *sql.DB, metrics *metrics.MicroserviceMetrics, statements map[string]*sql.Stmt) Repo {
	return &Layer{
		database: db,
		metrics:  metrics,
		stmt: statements,
	}
}

func (repo *Layer) GetByEmail(ctx context.Context, email *user.Email) (*user.User, error) {
	timeNow := time.Now()
	row := repo.stmt["getUser"].QueryRowContext(ctx, email.GetEmail())
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
	u := entity.User{}
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Address, &u.ImgUrl, &u.CardNumber, &u.IsVkUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsUserRelation
		}
		return &user.User{}, err
	}

	return entity.ConvertEntityUserIntoProtoUser(&u), nil
}

func (repo *Layer) DeleteByEmail(ctx context.Context, email *user.Email) error {
	timeNow := time.Now()
	row, err := repo.stmt["deleteUser"].ExecContext(ctx, email.GetEmail())
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.DELETE).Observe(float64(timeEnd.Milliseconds()))
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
	timeNowMetric := time.Now()

	row, err := repo.stmt["createUser"].ExecContext(ctx, u.GetName(), u.GetEmail(), functions.MaybeNullString(u.GetPhone()), 
	u.GetPassword(), functions.MaybeNullString(u.GetAddress()), functions.MaybeNullString(u.GetImgUrl()), u.GetIsVkUser(), timeNow, timeNow)
	timeEnd := time.Since(timeNowMetric)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.INSERT).Observe(float64(timeEnd.Milliseconds()))
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
	timeNowMetric := time.Now()
	row, err := repo.stmt["updateUser"].ExecContext(ctx,
		userData.GetName(), userData.GetEmail(), functions.MaybeNullString(userData.GetPhone()), userData.GetImgUrl(),
		userData.GetPassword(), functions.MaybeNullString(userData.GetCardNumber()), functions.MaybeNullString(userData.GetAddress()), timeNow, data.GetEmail())

	timeEnd := time.Since(timeNowMetric)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.UPDATE).Observe(float64(timeEnd.Milliseconds()))
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
	timeNow := time.Now()
	row := repo.stmt["getUnauthAddress"].QueryRowContext(ctx, id.GetUnauthId())
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.SELECT).Observe(float64(timeEnd.Milliseconds()))
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
	timeNow := time.Now()
	row, err := repo.stmt["updateUnauthAddress"].ExecContext(ctx, functions.MaybeNullString(data.GetAddress()), data.GetUnauthId())
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.UPDATE).Observe(float64(timeEnd.Milliseconds()))
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
	timeNow := time.Now()
	row, err := repo.stmt["createUnauthAddress"].ExecContext(ctx, data.GetUnauthId(), data.GetAddress())
	timeEnd := time.Since(timeNow)
	repo.metrics.DatabaseDuration.WithLabelValues(metrics.INSERT).Observe(float64(timeEnd.Milliseconds()))
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

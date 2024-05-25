package promocode

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
)

type Repo interface {
	GetPromocode(ctx context.Context, code string) (*entity.Promocode, error)
	WasPromocodeUsed(ctx context.Context, userId alias.UserId, codeId uint64) error
	WasRestPromocodeUsed(ctx context.Context, orderId alias.OrderId, codeId uint64) error
	SetPromocode(ctx context.Context, orderId alias.OrderId, codeId uint64) (uint64, error)
	GetPromocodeByOrder(ctx context.Context, orderId *alias.OrderId) (*entity.Promocode, error)
	DeletePromocode(ctx context.Context, orderId alias.OrderId) error
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(dbProps *sql.DB) Repo {
	return &RepoLayer{
		db: dbProps,
	}
}

// ПРОМОКОДЫ
func (repo *RepoLayer) GetPromocode(ctx context.Context, code string) (*entity.Promocode, error) {
	res := entity.PromocodeDB{}
	fmt.Println(code)
	err := repo.db.QueryRowContext(ctx,
		`SELECT id, date, sale, type, restaurant_id, sum FROM promocode WHERE code=$1`, code).Scan(&res.Id, &res.Date, &res.Sale, &res.Type, &res.Rest, &res.Sum)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsPromocodeRelation
		}
		return nil, err
	}
	res.Code = code
	return entity.ToPromocode(&res), nil
}

// может быть применен один раз, то есть он может быть тольок в одном заказе
func (repo *RepoLayer) WasPromocodeUsed(ctx context.Context, userId alias.UserId, codeId uint64) error {
	var res uint64
	err := repo.db.QueryRowContext(ctx,
		`SELECT count(*) FROM "order" WHERE user_id=$1 AND promocode_id=$2 AND status='delivered'`, uint64(userId), codeId).Scan(&res)
	if err != nil {
		fmt.Println(err, res)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if res > 1 {
		return myerrors.OncePromocode
	}
	return nil
}

func (repo *RepoLayer) WasRestPromocodeUsed(ctx context.Context, orderId alias.OrderId, codeId uint64) error {
	var res uint64
	err := repo.db.QueryRowContext(ctx,
		`SELECT count(*) FROM "order" WHERE id=$1 AND promocode_id=$2`, uint64(orderId), codeId).Scan(&res)
	if err != nil {
		fmt.Println(err, res)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if res > 1 {
		return myerrors.OncePromocode
	}
	return nil
}

func (repo *RepoLayer) SetPromocode(ctx context.Context, orderId alias.OrderId, codeId uint64) (uint64, error) {
	var res uint64
	err := repo.db.QueryRowContext(ctx,
		`UPDATE "order" SET promocode_id=$1 WHERE id=$2 RETURNING sum`, codeId, orderId).Scan(&res)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myerrors.SqlNoRowsPromocodeRelation
		}
		return 0, err
	}
	return res, nil
}

func (repo *RepoLayer) DeletePromocode(ctx context.Context, orderId alias.OrderId) error {
	res, err := repo.db.ExecContext(ctx,
		`UPDATE "order" SET promocode_id=NULL WHERE id=$1`, orderId)
	if err != nil {
		return err
	}
	countRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countRows == 0 {
		return myerrors.SqlNoRowsOrderRelation
	}
	return nil
}

func (repo *RepoLayer) GetPromocodeByOrder(ctx context.Context, orderId *alias.OrderId) (*entity.Promocode, error) {
	var i sql.NullInt64
	err := repo.db.QueryRowContext(ctx,
		`SELECT promocode_id FROM "order" WHERE id=$1`, orderId).Scan(&i)
	fmt.Println(err, i.Int64)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsPromocodeRelation
		}
		return nil, err
	}
	id := entity.Int(i)
	if id == 0 {
		fmt.Println("tut")
		return nil, myerrors.SqlNoRowsPromocodeRelation
	}
	res := entity.PromocodeDB{}
	err = repo.db.QueryRowContext(ctx,
		`SELECT id, code, date, sale, type, restaurant_id, sum FROM promocode WHERE id=$1`, id).Scan(&res.Id, &res.Code, &res.Date, &res.Sale, &res.Type, &res.Rest, &res.Sum)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.SqlNoRowsPromocodeRelation
		}
		return nil, err
	}
	return entity.ToPromocode(&res), nil
}

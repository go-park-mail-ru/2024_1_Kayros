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
	stmt map[string]*sql.Stmt
}

func NewRepoLayer(dbProps *sql.DB, statements map[string]*sql.Stmt) Repo {
	return &RepoLayer{
		db: dbProps,
		stmt: statements,
	}
}

// ПРОМОКОДЫ
func (repo *RepoLayer) GetPromocode(ctx context.Context, code string) (*entity.Promocode, error) {
	res := entity.PromocodeDB{}
	fmt.Println(code)
	err := repo.stmt["getPromocodeByCode"].QueryRowContext(ctx, code).Scan(&res.Id, &res.Date, &res.Sale, &res.Type, &res.Rest, &res.Sum)
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
	err := repo.stmt["wasPromocodeUsed"].QueryRowContext(ctx, uint64(userId), codeId).Scan(&res)
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
	err := repo.stmt["wasRestPromocodeUsed"].QueryRowContext(ctx, uint64(orderId), codeId).Scan(&res)
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
	err := repo.stmt["setPromocode"].QueryRowContext(ctx, codeId, orderId).Scan(&res)
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
	res, err := repo.stmt["deletePromocode"].ExecContext(ctx, orderId)
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
	err := repo.stmt["getPromocodeIdFromOrder"].QueryRowContext(ctx, orderId).Scan(&i)
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
	err = repo.stmt["getPromocodeById"].QueryRowContext(ctx, id).Scan(&res.Id, &res.Code, &res.Date, &res.Sale, &res.Type, &res.Rest, &res.Sum)
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

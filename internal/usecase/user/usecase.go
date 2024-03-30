package user

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetById(context.Context, alias.UserId) (*entity.User, error)
	GetByEmail(context.Context, string) (*entity.User, error)
	GetByPhone(context.Context, string) (*entity.User, error)

	DeleteById(context.Context, alias.UserId) (bool, error)
	DeleteByEmail(context.Context, string) (bool, error)
	DeleteByPhone(context.Context, string) (bool, error)

	Create(context.Context, *entity.User) (*entity.User, error)
	Update(context.Context, *entity.User) (*entity.User, error)

	IsExistById(context.Context, alias.UserId) (bool, error)
	IsExistByEmail(context.Context, string) (bool, error)

	CheckPassword(context.Context, alias.UserId, string) (bool, error)
}

type UsecaseLayer struct {
	repo user.Repo
}

func NewUsecaseLayer(repoUser user.Repo) Usecase {
	return &UsecaseLayer{
		repo: repoUser,
	}
}

func (uc *UsecaseLayer) GetById(ctx context.Context, id alias.UserId) (*entity.User, error) {
	u, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UsecaseLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UsecaseLayer) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	u, err := uc.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UsecaseLayer) DeleteById(ctx context.Context, id alias.UserId) (bool, error) {
	wasDeleted, err := uc.DeleteById(ctx, id)
	return wasDeleted, err
}

func (uc *UsecaseLayer) DeleteByEmail(ctx context.Context, email string) (bool, error) {
	wasDeleted, err := uc.DeleteByEmail(ctx, email)
	return wasDeleted, err
}

func (uc *UsecaseLayer) DeleteByPhone(ctx context.Context, phone string) (bool, error) {
	wasDeleted, err := uc.DeleteByPhone(ctx, phone)
	return wasDeleted, err
}

func (uc *UsecaseLayer) Create(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	u, err := uc.repo.Create(ctx, uProps)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UsecaseLayer) Update(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	u, err := uc.repo.Update(ctx, uProps)
	if err != nil {
		return nil, err
	}
	return u, err
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (uc *UsecaseLayer) CheckPassword(ctx context.Context, id alias.UserId, password string) (bool, error) {
	isEqual, err := uc.CheckPassword(ctx, id, password)
	return isEqual, err
}

func (uc *UsecaseLayer) IsExistById(ctx context.Context, id alias.UserId) (bool, error) {
	isExist, err := uc.IsExistById(ctx, id)
	return isExist, err
}

func (uc *UsecaseLayer) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	isExist, err := uc.IsExistByEmail(ctx, email)
	return isExist, err
}

package user

import (
	"context"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/alias"
)

type Usecase interface {
	GetById(ctx context.Context, userId alias.UserId) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	DeleteById(ctx context.Context, userId alias.UserId) (bool, error)
	DeleteByEmail(ctx context.Context, email string) (bool, error)

	IsExistById(ctx context.Context, userId alias.UserId) (bool, error)
	IsExistByEmail(ctx context.Context, email string) (bool, error)

	Create(ctx context.Context, uProps *entity.User) (*entity.User, error)
	Update(ctx context.Context, uProps *entity.User) (*entity.User, error)

	CheckPassword(ctx context.Context, email string, password string) (bool, error)
}

type UsecaseLayer struct {
	repoUser user.Repo
}

func NewUsecaseLayer(repoUserProps user.Repo) Usecase {
	return &UsecaseLayer{
		repoUser: repoUserProps,
	}
}

func (uc *UsecaseLayer) GetById(ctx context.Context, userId alias.UserId) (*entity.User, error) {
	return uc.repoUser.GetById(ctx, userId)
}

func (uc *UsecaseLayer) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return uc.repoUser.GetByEmail(ctx, email)
}

func (uc *UsecaseLayer) DeleteById(ctx context.Context, userId alias.UserId) (bool, error) {
	return uc.repoUser.DeleteById(ctx, userId)
}

func (uc *UsecaseLayer) DeleteByEmail(ctx context.Context, email string) (bool, error) {
	return uc.repoUser.DeleteByEmail(ctx, email)
}

func (uc *UsecaseLayer) IsExistById(ctx context.Context, userId alias.UserId) (bool, error) {
	return uc.repoUser.IsExistById(ctx, userId)
}

func (uc *UsecaseLayer) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	return uc.repoUser.IsExistByEmail(ctx, email)
}

func (uc *UsecaseLayer) Create(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	return uc.repoUser.Create(ctx, uProps)
}

func (uc *UsecaseLayer) Update(ctx context.Context, uProps *entity.User) (*entity.User, error) {
	return uc.repoUser.Update(ctx, uProps)
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (uc *UsecaseLayer) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
	return uc.repoUser.CheckPassword(ctx, email, password)
}

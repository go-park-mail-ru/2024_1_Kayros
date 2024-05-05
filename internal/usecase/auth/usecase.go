package auth

import (
	"bytes"
	"context"
	"errors"

	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type Usecase interface {
	SignUpUser(ctx context.Context, email string, unauthId string, signupData *entity.User) (*entity.User, error)
	SignInUser(ctx context.Context, email string, unauthId string, password string) (*entity.User, error)
}

type UsecaseLayer struct {
	repoUser user.Repo
}

func NewUsecaseLayer(repoUserProps user.Repo) Usecase {
	return &UsecaseLayer{
		repoUser: repoUserProps,
	}
}

func (uc *UsecaseLayer) SignUpUser(ctx context.Context, email string, unauthId string, signupData *entity.User) (*entity.User, error) {
	isExist, err := uc.isExistByEmail(ctx, email)
	if err != nil {
		// we can skip error `myerrors.SqlNoRowsUserRelation`, because user must not to be
		if !errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return nil, err
		}
	}
	if isExist {
		return nil, myerrors.UserAlreadyExist
	}

	address, err := uc.repoUser.GetAddressByUnauthId(ctx, unauthId)
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
		return nil, err
	}
	signupData.Address = address

	// we do copy for clean function
	uCopy := entity.Copy(signupData)
	salt, err := functions.GenerateNewSalt()
	if err != nil {
		return nil, err
	}
	hashPassword := functions.HashData(salt, signupData.Password)
	uCopy.Password = string(hashPassword)

	err = uc.repoUser.Create(ctx, uCopy)
	if err != nil {
		return nil, err
	}

	uDB, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return uDB, nil
}

func (uc *UsecaseLayer) SignInUser(ctx context.Context, email string, unauthId string, password string) (*entity.User, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	isEqual, err := uc.checkPassword(ctx, email, password)
	if err != nil {
		return nil, err
	}
	if !isEqual {
		return nil, myerrors.BadAuthPassword
	}

	address, err := uc.repoUser.GetAddressByUnauthId(ctx, unauthId)
	if err != nil && !errors.Is(err, myerrors.SqlNoRowsUnauthAddressRelation) {
		return nil, err
	}
	if address != "" && u.Address == "" {
		uDataChange := entity.Copy(u)
		uDataChange.Address = address
		err = uc.repoUser.Update(ctx, uDataChange, email)
		if err != nil {
			return nil, err
		}
		u.Address = address
	}
	return u, nil
}

func (uc *UsecaseLayer) isExistByEmail(ctx context.Context, email string) (bool, error) {
	_, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, myerrors.SqlNoRowsUserRelation) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// checkPassword - method used to check password with password saved in database
func (uc *UsecaseLayer) checkPassword(ctx context.Context, email string, password string) (bool, error) {
	u, err := uc.repoUser.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	uPasswordBytes := []byte(u.Password)

	salt := make([]byte, 8)
	copy(salt, uPasswordBytes[0:8])
	hashPassword := functions.HashData(salt, password)
	return bytes.Equal(uPasswordBytes, hashPassword), nil
}

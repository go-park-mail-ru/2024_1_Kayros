package user

import (
	"encoding/json"
	"net/http"

	"2024_1_kayros/internal/entity/dto"
	"2024_1_kayros/internal/repository/user"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
)

type UserUsecaseInterface interface {
	GetData(w http.ResponseWriter, r *http.Request) http.ResponseWriter
}

type UserUsecase struct {
	repoUser user.UserRepositoryInterface
}

func NewUserUsecase(repoUser user.UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		repoUser: repoUser,
	}
}

func (uc *UserUsecase) GetData(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	// если пришел неавторизованный пользователь, возвращаем 401
	authKey := r.Context().Value("authKey")
	if authKey == nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
		return w
	}

	u, err := uc.repoUser.GetByEmail(authKey.(string))
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusUnauthorized)
	}

	// Собираем ответ
	uDTO := &dto.UserDTO{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		ImgUrl:   u.ImgUrl,
		Password: u.Password,
	}
	jsonResponse, err := json.Marshal(uDTO)
	if err != nil {
		w = functions.ErrorResponse(w, myerrors.InternalServerError, http.StatusInternalServerError)
		return w
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w
	}

	w.WriteHeader(http.StatusOK)
	return w
}

package user

import (
	"context"
	"io"
	"mime/multipart"

	protouser "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/entity"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/internal/utils/props"

	"google.golang.org/grpc/codes"
)

type Usecase interface {
	UserAddress(ctx context.Context, email, unauthId string) (string, error)
	GetData(ctx context.Context, email string) (*entity.User, error)
	UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error)
	UpdateAddress(ctx context.Context, email, unauthId, address string) error
	SetNewPassword(ctx context.Context, email, password, newPassword string) error
}

type UsecaseLayer struct {
	userClient protouser.UserManagerClient
}

func NewUsecaseLayer(userClientProps protouser.UserManagerClient) Usecase {
	return &UsecaseLayer{
		userClient: userClientProps,
	}
}

// UserAddress - method returns user address by unauthId or email, but unauthId in priority
func (uc *UsecaseLayer) UserAddress(ctx context.Context, email, unauthId string) (string, error) {
	address := ""
	u, err := uc.userClient.GetData(ctx, &protouser.Email{Email: email})
	if err != nil && !grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
		return "", err
	}
	unauthAddress, err := uc.userClient.GetAddressByUnauthId(ctx, &protouser.UnauthId{UnauthId: unauthId})
	if err != nil && !grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUnauthAddressRelation) {
		return "", err
	}
	if unauthAddress.GetAddress() != "" {
		address = unauthAddress.GetAddress()
	} else if !entity.ProtoUserIsNIL(u) && u.Address != "" {
		address = u.Address
	}
	return address, nil
}

// GetData - method calls repo method to receive user data.
func (uc *UsecaseLayer) GetData(ctx context.Context, email string) (*entity.User, error) {
	u, err := uc.userClient.GetData(ctx, &protouser.Email{Email: email})
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return &entity.User{}, myerrors.SqlNoRowsUserRelation
		}
		return &entity.User{}, err
	}
	return entity.ConvertProtoUserIntoEntityUser(u), nil
}

// UpdateData - method used to update user info. It accepts non-password fields.
// To update password use method SetNewUserPassword.
func (uc *UsecaseLayer) UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error) {
	var fileData []byte
	var err error
	if data.File != nil {
		fileData, err = MultipartFileToBytes(data.File)
		if err != nil {
			return &entity.User{}, err
		}
	}
	var fileName string
	var fileSize int64
	if data.Handler != nil {
		fileName = data.Handler.Filename
		fileSize = data.Handler.Size
	}
	dataProps := &protouser.UpdateUserData{
		UpdateInfo: entity.ConvertEntityUserIntoProtoUser(data.UserPropsUpdate),
		Email:      data.Email,
		FileData:   fileData,
		FileName:   fileName,
		FileSize:   fileSize,
	}
	u, err := uc.userClient.UpdateData(ctx, dataProps)
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return &entity.User{}, myerrors.SqlNoRowsUserRelation
		}
		if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUserRelation) {
			return &entity.User{}, myerrors.SqlNoRowsUserRelation
		}
		if grpcerr.Is(err, codes.InvalidArgument, myerrors.WrongFileExtension) {
			return &entity.User{}, myerrors.WrongFileExtension
		}
		if grpcerr.Is(err, codes.AlreadyExists, myerrors.UserAlreadyExist) {
			return &entity.User{}, myerrors.UserAlreadyExist
		}
		return &entity.User{}, err
	}
	return entity.ConvertProtoUserIntoEntityUser(u), nil
}

// SetNewPassword - method used to set new password.
func (uc *UsecaseLayer) SetNewPassword(ctx context.Context, email, password, newPassword string) error {
	data := &protouser.PasswordsChange{
		Email:       email,
		Password:    password,
		NewPassword: newPassword,
	}
	_, err := uc.userClient.SetNewPassword(ctx, data)
	if err != nil {
		if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
			return myerrors.SqlNoRowsUserRelation
		}
		if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUserRelation) {
			return myerrors.SqlNoRowsUserRelation
		}
		if grpcerr.Is(err, codes.InvalidArgument, myerrors.IncorrectCurrentPassword) {
			return myerrors.IncorrectCurrentPassword
		}
		if grpcerr.Is(err, codes.InvalidArgument, myerrors.NewPassword) {
			return myerrors.NewPassword
		}
		return err
	}
	return nil
}

// UpdateAddress - method updates only address.
func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, email, unauthId, address string) error {
	if email != "" {
		data := &protouser.AddressData{
			Email:   email,
			Address: address,
		}
		_, err := uc.userClient.UpdateAddress(ctx, data)
		if err != nil {
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
				return myerrors.SqlNoRowsUserRelation
			}
			if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUserRelation) {
				return myerrors.SqlNoRowsUserRelation
			}
			return err
		}
	}
	if unauthId != "" {
		data := &protouser.AddressDataUnauth{
			UnauthId: unauthId,
			Address:  address,
		}
		_, err := uc.userClient.UpdateAddressByUnauthId(ctx, data)
		if err != nil {
			if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUnauthAddressRelation) {
				return myerrors.SqlNoRowsUnauthAddressRelation
			}
			return err
		}
	}
	return nil
}

func MultipartFileToBytes(file multipart.File) ([]byte, error) {
	defer file.Close()

	// Читаем данные из файла
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

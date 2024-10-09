package user

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	protouser "2024_1_kayros/gen/go/user"
	"2024_1_kayros/internal/delivery/metrics"
	"2024_1_kayros/internal/entity"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/myerrors/grpcerr"
	"2024_1_kayros/internal/utils/props"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase interface {
	UserAddress(ctx context.Context, email, unauthId, isUserAddress string) (string, error)
	GetData(ctx context.Context, email string) (*entity.User, error)
	UpdateData(ctx context.Context, data *props.UpdateUserDataProps) (*entity.User, error)
	UpdateAddress(ctx context.Context, email, unauthId, addressm, isUserAddress string) error
	UpdateUnauthAddress(ctx context.Context, address string, unauthId string) error
	SetNewPassword(ctx context.Context, email, password, newPassword string) error
}

type UsecaseLayer struct {
	userClient protouser.UserManagerClient
	metrics    *metrics.Metrics
}

func NewUsecaseLayer(userClientProps protouser.UserManagerClient, metrics *metrics.Metrics) Usecase {
	return &UsecaseLayer{
		userClient: userClientProps,
		metrics:    metrics,
	}
}

// UserAddress - method returns user address by unauthId or email, but unauthId in priority
func (uc *UsecaseLayer) UserAddress(ctx context.Context, email, unauthId, isUserAddress string) (string, error) {
	address := ""
	timeNow := time.Now()
	u, err := uc.userClient.GetData(ctx, &protouser.Email{Email: email})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil && !grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
		}
		return "", err
	}
	if isUserAddress == "true" {
		return u.Address, nil
	}
	if unauthId != "" {
		timeNow = time.Now()
		unauthAddress, err := uc.userClient.GetAddressByUnauthId(ctx, &protouser.UnauthId{UnauthId: unauthId})
		msRequestTimeout = time.Since(timeNow)
		uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
		if err != nil && !grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUnauthAddressRelation) {
			grpcStatus, ok := status.FromError(err)
			if !ok {
				uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
			}
			return "", err
		}
		if unauthAddress != nil && unauthAddress.GetAddress() != "" {
			address = unauthAddress.GetAddress()
		}
		return address, nil
	}
	if !entity.ProtoUserIsNIL(u) && u.Address != "" {
		address = u.Address
	}
	return address, nil
}

func (uc *UsecaseLayer) UpdateUnauthAddress(ctx context.Context, address string, unauthId string) error {
	timeNow := time.Now()
	_, err := uc.userClient.UpdateAddressByUnauthId(ctx, &protouser.AddressDataUnauth{UnauthId: unauthId, Address: address})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
		}
		return err
	}
	return nil
}

// GetData - method calls repo method to receive user data.
func (uc *UsecaseLayer) GetData(ctx context.Context, email string) (*entity.User, error) {
	timeNow := time.Now()
	u, err := uc.userClient.GetData(ctx, &protouser.Email{Email: email})
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
		}
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
	timeNow := time.Now()
	u, err := uc.userClient.UpdateData(ctx, dataProps)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
		}
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
	timeNow := time.Now()
	_, err := uc.userClient.SetNewPassword(ctx, data)
	msRequestTimeout := time.Since(timeNow)
	uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if !ok {
			uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
		}
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
func (uc *UsecaseLayer) UpdateAddress(ctx context.Context, email, unauthId, address, isUserAddress string) error {
	if email != "" {
		data := &protouser.AddressData{
			Email:   email,
			Address: address,
		}
		timeNow := time.Now()
		_, err := uc.userClient.UpdateAddress(ctx, data)
		msRequestTimeout := time.Since(timeNow)
		uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
		if err != nil {
			grpcStatus, ok := status.FromError(err)
			if !ok {
				uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
			}
			if grpcerr.Is(err, codes.NotFound, myerrors.SqlNoRowsUserRelation) {
				return myerrors.SqlNoRowsUserRelation
			}
			if grpcerr.Is(err, codes.Internal, myerrors.SqlNoRowsUserRelation) {
				return myerrors.SqlNoRowsUserRelation
			}
			return err
		}
		if isUserAddress == "true" {
			return nil
		}
	}
	if unauthId != "" {
		data := &protouser.AddressDataUnauth{
			UnauthId: unauthId,
			Address:  address,
		}
		timeNow := time.Now()
		_, err := uc.userClient.UpdateAddressByUnauthId(ctx, data)
		msRequestTimeout := time.Since(timeNow)
		uc.metrics.MicroserviceTimeout.WithLabelValues(cnst.UserMicroservice).Observe(float64(msRequestTimeout.Milliseconds()))
		if err != nil {
			grpcStatus, ok := status.FromError(err)
			if !ok {
				uc.metrics.MicroserviceErrors.WithLabelValues(cnst.UserMicroservice, grpcStatus.String()).Inc()
			}
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

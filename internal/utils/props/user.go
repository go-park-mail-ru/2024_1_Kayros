package props

import (
	"mime/multipart"

	"2024_1_kayros/internal/entity"
)

// UpdateUserDataProps - props used in UpdateUserData (usecase) | User
type UpdateUserDataProps struct {
	Email           string
	File            multipart.File
	Handler         *multipart.FileHeader
	UserPropsUpdate *entity.User
}

func GetUpdateUserDataProps(email string, file multipart.File, handler *multipart.FileHeader, u *entity.User) *UpdateUserDataProps {
	return &UpdateUserDataProps{
		Email:           email,
		File:            file,
		Handler:         handler,
		UserPropsUpdate: u,
	}
}

// SetNewUserPasswordProps - props used in SetNewUserPassword (usecase) !!!
type SetNewUserPasswordProps struct {
	Password    string
	PasswordNew string
}

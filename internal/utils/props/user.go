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

// SetNewUserPasswordProps - props used in
type SetNewUserPasswordProps struct {
	Password    string
	PasswordNew string
}

package myerrors

import "errors"

// Response errors (external) RU
var (
	InternalServerRu     = errors.New("Ошибка сервера")
	BadAuthCredentialsRu = errors.New("Неверный логин или пароль")

	UnauthorizedRu = errors.New("Вы не авторизованы")
	RegisteredRu   = errors.New("Вы уже зарегистрированы")
	AuthorizedRu   = errors.New("Вы уже авторизованы")

	BadCredentialsRu           = errors.New("Предоставлены некорректные данные")
	UserAlreadyExistRu         = errors.New("Пользователь с таким логином уже зарегистрирован")
	BigSizeFileRu              = errors.New("Превышен максимальный размер файла")
	NotFoundRu                 = errors.New("Данные не найдены")
	WrongFileExtensionRu       = errors.New("Недопустимый формат фотографии")
	NewPasswordRu              = errors.New("Новый пароль совпадает со старым")
	IncorrectCurrentPasswordRu = errors.New("Неверно указан текущий пароль")
)

// Response errors (external) En
var (
	InternalServerEn     = errors.New("Server error")
	BadAuthCredentialsEn = errors.New("Invalid login or password")

	UnauthorizedEn = errors.New("You are not authorized")
	RegisteredEn   = errors.New("You are already registered")
	AuthorizedEn   = errors.New("You are already authorized")

	BadCredentialsEn           = errors.New("Incorrect data provided")
	UserAlreadyExistEn         = errors.New("User with this login already exists")
	BigSizeFileEn              = errors.New("Maximum file size exceeded")
	NotFoundEn                 = errors.New("Data not found")
	WrongFileExtensionEn       = errors.New("Invalid image format")
	NewPasswordEn              = errors.New("The new password is the same as the old one")
	IncorrectCurrentPasswordEn = errors.New("The current password is incorrect")
	NameEn                     = errors.New("Incorrect name specified")
	PhoneEn                    = errors.New("Incorrect phone number specified")
	EmailEn                    = errors.New("Incorrect email specified")
	AddressEn                  = errors.New("Incorrect address specified")
)

// Internal errors
var (
	BigSizeFile              = errors.New("maximum file size exceeded")
	CtxRequestId             = errors.New("request_id was not passed in the context")
	CtxEmail                 = errors.New("email was not passed in the context")
	WrongFileExtension       = errors.New("invalid image format")
	NewPassword              = errors.New("the new password is the same as the old one")
	IncorrectCurrentPassword = errors.New("the current password is incorrect")
	UserAlreadyExist         = errors.New("user with this login already exists")
	BadAuthPassword          = errors.New("invalid password")
	// Database
	SqlNoRowsUserRelation = errors.New("no such record exists for \"user\"")
	RedisNoData           = errors.New("no such record exists in Redis")
)

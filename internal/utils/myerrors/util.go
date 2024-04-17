package myerrors

import "errors"

// Response errors (external) RU
var (
	InternalServerRu     = errors.New("Ошибка сервера")
	BadAuthCredentialsRu = errors.New("Неверный логин или пароль")
	UnauthorizedRu       = errors.New("Вы не зарегистрированы")
	BadCredentialsRu     = errors.New("Предоставлены некорректные данные")
	UserAlreadyExistRu   = errors.New("Пользователь с таким логином уже зарегистрирован")
	BigSizeFileRu        = errors.New("Превышен максимальный размер файла")
	NotFoundRu           = errors.New("Данные не найдены")
	WrongFileExtensionRu = errors.New("Недопустимый формат фотографии")
	NewPasswordRu        = errors.New("Новый пароль совпадает со старым")
	PasswordRu           = errors.New("Неверно указан текущий пароль")
	NameRu               = errors.New("Неверно указано имя")
	PhoneRu              = errors.New("Неверно указан номер телефона ")
	EmailRu              = errors.New("Неверно указана почта")
	AddressRu            = errors.New("Неверно указан адрес")
)

// Response errors (external) En
var (
	InternalServerEn     = errors.New("Server error")
	BadAuthCredentialsEn = errors.New("Invalid login or password")
	UnauthorizedEn       = errors.New("You are not registered")
	BadCredentialsEn     = errors.New("Incorrect data provided")
	UserAlreadyExistEn   = errors.New("User with this login already exists")
	BigSizeFileEn        = errors.New("Maximum file size exceeded")
	NotFoundEn           = errors.New("Data not found")
	WrongFileExtensionEn = errors.New("Invalid image format")
	NewPasswordEn        = errors.New("The new password is the same as the old one")
	PasswordEn           = errors.New("The current password is incorrect")
	NameEn               = errors.New("Incorrect name specified")
	PhoneEn              = errors.New("Incorrect phone number specified")
	EmailEn              = errors.New("Incorrect email specified")
	AddressEn            = errors.New("Incorrect address specified")
)

// Internal errors
var (
	BigSizeFile        = errors.New("maximum file size exceeded")
	CtxRequestId       = errors.New("request_id was not passed in the context")
	CtxEmail           = errors.New("email was not passed in the context")
	WrongFileExtension = errors.New("invalid image format")
	NewPassword        = errors.New("the new password is the same as the old one")
	Password           = errors.New("the current password is incorrect")
	Name               = errors.New("invalid name")
	Phone              = errors.New("invalid phone")
	Email              = errors.New("invalid email")
	Address            = errors.New("the address exceeds the permissible size")
	// Database
	SqlNoRowsUserRelation = errors.New("no such record exists for \"user\"")
	RedisNoData           = errors.New("no such record exists in Redis")
)

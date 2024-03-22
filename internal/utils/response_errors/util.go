package response_errors

import "errors"

var (
	UnexpectedServerError   = errors.New("Ошибка сервера")
	BadPermissionError      = errors.New("Не хватает прав для доступа")
	BadAuthCredentialsError = errors.New("Неверный логин или пароль")
	BadRegCredentialsError  = errors.New("Некорректные данные")
	UserAlreadyExistError   = errors.New("Пользователь с таким логином уже зарегистрирован")
)

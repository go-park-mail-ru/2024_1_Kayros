package myerrors

import "errors"

var (
	UnexpectedServerError   = errors.New("Ошибка сервера")
	HashedPasswordError     = errors.New("Не удалось захешировать пароль пользователя")
	BadPermissionError      = errors.New("Не хватает прав для доступа")
	BadAuthCredentialsError = errors.New("Неверный логин или пароль")
	BadRegCredentialsError  = errors.New("Некорректные данные")
	UserAlreadyExistError   = errors.New("Пользователь с таким логином уже зарегистрирован")
)

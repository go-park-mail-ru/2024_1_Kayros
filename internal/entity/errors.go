package entity

import "errors"

var (
	UnexpectedServerError = errors.New("Ошибка сервера")
	BadPermission         = errors.New("Не хватает прав для доступа")
	BadAuthCredentials    = errors.New("Неверный логин или пароль")
	BadRegCredentials     = errors.New("Некорректные данные")
	// UserAlreadyExist = errors.New("Пользователь с таким логином уже зарегистрирован")
)

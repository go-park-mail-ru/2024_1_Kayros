package myerrors

var (
	IntServerError          = "Ошибка сервера"
	HashedPasswordError     = "Не удалось захешировать пароль пользователя"
	BadPermissionError      = "Не хватает прав для доступа"
	BadAuthCredentialsError = "Неверный логин или пароль"
	BadCredentialsError     = "Предоставлены некорректные данные"
	UserAlreadyExistError   = "Пользователь с таким логином уже зарегистрирован"
)

package myerrors

const (
	InternalServerError     = "Ошибка сервера"
	HashedPasswordError     = "Не удалось захешировать пароль пользователя"
	UnauthorizedError       = "Не хватает прав для доступа"
	BadAuthCredentialsError = "Неверный логин или пароль"
	BadCredentialsError     = "Предоставлены некорректные данные"
	UserAlreadyExistError   = "Пользователь с таким логином уже зарегистрирован"
	BigSizeFileError        = "Превышен максимальный размер файла"
	NotFoundError           = "Данные не найдены"
)

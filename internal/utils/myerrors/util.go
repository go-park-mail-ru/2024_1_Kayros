package myerrors

// Ошибки сетевых ответов
const (
	InternalServerError     = "Ошибка сервера"
	UnauthorizedError       = "Вы не зарегистрированы"
	BadAuthCredentialsError = "Неверный логин или пароль"
	BadCredentialsError     = "Предоставлены некорректные данные"
	UserAlreadyExistError   = "Пользователь с таким логином уже зарегистрирован"
	BigSizeFileError        = "Превышен максимальный размер файла"
	NotFoundError           = "Данные не найдены"
	EqualPasswordsError     = "Новый пароль совпадает со старым"
	WrongPasswordError      = "Не верно введен старый пароль"
)

// Внутренние ошибки
const (
	HashedPasswordError = "Не удалось захешировать пароль пользователя"
)

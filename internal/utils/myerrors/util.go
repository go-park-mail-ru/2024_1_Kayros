package myerrors

import "errors"

// Response errors (external) RUSSIA
var (
	InternalServerRu     = errors.New("Ошибка сервера")
	BadAuthCredentialsRu = errors.New("Неверный логин или пароль")
	BadCredentialsRu     = errors.New("Предоставлены некорректные данные")
	UserAlreadyExistRu   = errors.New("Пользователь с таким логином уже зарегистрирован")
	NotFoundRu           = errors.New("Данные не найдены")

	UnauthorizedRu   = errors.New("Вы не авторизованы")
	RegisteredRu     = errors.New("Вы уже зарегистрированы")
	AuthorizedRu     = errors.New("Вы уже авторизованы")
	SignOutAlreadyRu = errors.New("Вы уже вышли из аккаунта")

	BigSizeFileRu              = errors.New("Превышен максимальный размер файла")
	WrongFileExtensionRu       = errors.New("Недопустимый формат фотографии")
	NewPasswordRu              = errors.New("Новый пароль совпадает со старым")
	IncorrectCurrentPasswordRu = errors.New("Неверно указан текущий пароль")

	QuizAddRu = errors.New("Произошла ошибка. Пожалуйста, еще раз проголосуйте")

	NoBasketRu         = errors.New("У Вас нет корзины")
	AlreadyPayedRu     = errors.New("Заказ уже оплачен")
	NoAddFoodToOrderRu = errors.New("Не удалось добавить блюдо в заказ, попробуйте еще раз")
	FailCleanBasketRu  = errors.New("Не удалось очистить корзину")
	NoDeleteFoodRu     = errors.New("Не удалось убрать блюдо из заказа, попробуйте еще раз")
)

// Response errors (external) ENGLISH
var (
	InternalServerEn     = errors.New("Server error")
	BadAuthCredentialsEn = errors.New("Invalid login or password")
	BadCredentialsEn     = errors.New("Incorrect data provided")
	UserAlreadyExistEn   = errors.New("User with this login already exists")
	NotFoundEn           = errors.New("Data not found")

	UnauthorizedEn = errors.New("You are not authorized")
	RegisteredEn   = errors.New("You are already registered")
	AuthorizedEn   = errors.New("You are already authorized")

	BigSizeFileEn              = errors.New("Maximum file size exceeded")
	WrongFileExtensionEn       = errors.New("Invalid image format")
	NewPasswordEn              = errors.New("The new password is the same as the old one")
	IncorrectCurrentPasswordEn = errors.New("The current password is incorrect")
)

// Internal errors
var (
	UserAlreadyExist = errors.New("user with this login already exists")
	CtxRequestId     = errors.New("request_id was not passed in the context")
	CtxEmail         = errors.New("email was not passed in the context")

	BigSizeFile        = errors.New("maximum file size exceeded")
	WrongFileExtension = errors.New("invalid image format")

	NewPassword              = errors.New("the new password is the same as the old one")
	IncorrectCurrentPassword = errors.New("the current password is incorrect")
	BadAuthPassword          = errors.New("invalid password")

	BasketCreate     = errors.New("can't create basket")
	OrderAddFood     = errors.New("food was not added to the order")
	OrderSum         = errors.New("sum of the order is null")
	FailCleanBasket  = errors.New("can't clean basket")
	InvalidAddressEn = errors.New("invalid length of address")

	QuizAdd      = errors.New("answer was not added to the quiz")
	NoBasket     = errors.New("basket doesn't exist")
	AlreadyPayed = errors.New("order has already been paid")

	// Database
	SqlNoRowsUserRelation       = errors.New("no such record exists for \"user\"")
	SqlNoRowsRestaurantRelation = errors.New("no such record exists for restaurant")
	SqlNoRowsFoodRelation       = errors.New("no such record exists for food")
	SqlNoRowsOrderRelation      = errors.New("no such record exists for \"order\"")
	SqlNoRowsFoodOrderRelation  = errors.New("no such record exists for food_order")
	SqlNoRowsQuizRelation       = errors.New("no such record exists for quiz")
	RedisNoData                 = errors.New("no such record exists in Redis")
)

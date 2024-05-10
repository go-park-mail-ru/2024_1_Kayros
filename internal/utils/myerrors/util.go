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

	NoCommentsRu = errors.New(
		"У ресторана пока нет отзывов",
	)
	SuccessCleanRu      = errors.New("Корзина очищена")
	NoOrdersRu          = errors.New("Нет заказов")
	NoBasketRu          = errors.New("У Вас нет корзины")
	AlreadyPayedRu      = errors.New("Заказ уже оплачен")
	NoAddFoodToOrderRu  = errors.New("Не удалось добавить блюдо в заказ, попробуйте еще раз")
	FailCleanBasketRu   = errors.New("Не удалось очистить корзину")
	FailCreateCommentRu = errors.New("Не удалось добавить отзыв")
	NoDeleteFoodRu      = errors.New("Не удалось убрать блюдо из заказа, попробуйте еще раз")
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
	UserAlreadyExist = errors.New("Error: user with this login already exists")
	CtxRequestId     = errors.New("Error: request_id was not passed in the context")
	CtxEmail         = errors.New("Error: email was not passed in the context")

	BigSizeFile        = errors.New("Error: maximum file size exceeded")
	WrongFileExtension = errors.New("Error: invalid image format")

	NewPassword              = errors.New("Error: the new password is the same as the old one")
	IncorrectCurrentPassword = errors.New("Error: the current password is incorrect")
	BadAuthPassword          = errors.New("Error: invalid password")

	OrderAddFood     = errors.New("Error: food was not added to the order")
	OrderSum         = errors.New("Error: sum of the order is null")
	FailCleanBasket  = errors.New("Error: can't clean basket")
	BasketCreate     = errors.New("Error: can't create basket")
	InvalidAddressEn = errors.New("Error: invalid length of address")

	NoComments   = errors.New("Error: restaurant doesn't have comments")
	QuizAdd      = errors.New("Error: answer was not added to the quiz")
	NoBasket     = errors.New("Error: basket doesn't exist")
	AlreadyPayed = errors.New("Error: order has already been paid")

	// Database
	SqlNoRowsUserRelation          = errors.New("Error: no such record exists for \"user\"")
	SqlNoRowsUnauthAddressRelation = errors.New("Error: no such record exists for unauth_address")
	SqlNoRowsRestaurantRelation    = errors.New("Error: no such record exists for restaurant")
	SqlNoRowsFoodRelation          = errors.New("Error: no such record exists for food")
	SqlNoRowsOrderRelation         = errors.New("Error: no such record exists for \"order\"")
	SqlNoRowsCommentRelation       = errors.New("Error: no such record exists for \"comment\"")
	SqlNoRowsFoodOrderRelation     = errors.New("Error: no such record exists for food_order")
	SqlNoRowsQuizRelation          = errors.New("Error: no such record exists for quiz")
	RedisNoData                    = errors.New("Error: no such record exists in Redis")
	NullData                       = errors.New("Error: selected null data from table")
)

package myerrors

import "errors"

// Response errors (external) RUSSIA
var (
	InternalServerRu              = errors.New("Ошибка сервера")
	BadAuthCredentialsRu          = errors.New("Неверный логин или пароль")
	BadCredentialsRu              = errors.New("Предоставлены некорректные данные")
	UserAlreadyExistRu            = errors.New("Пользователь с таким логином уже зарегистрирован")
	NotFoundRu                    = errors.New("Данные не найдены")
	BadRequestGetAddress          = errors.New("Необходимо авторизоваться для получения адреса")
	BadRequestUpdateAddress       = errors.New("Необходимо авторизоваться для изменения адреса")
	BadRequestUpdateUnauthAddress = errors.New("Невозможно изменить адрес")

	UnauthorizedRu   = errors.New("Вы не авторизованы")
	RegisteredRu     = errors.New("Вы уже зарегистрированы")
	AuthorizedRu     = errors.New("Вы уже авторизованы")
	SignOutAlreadyRu = errors.New("Вы уже вышли из аккаунта")

	BigSizeFileRu              = errors.New("Превышен максимальный размер файла")
	WrongFileExtensionRu       = errors.New("Недопустимый формат фотографии")
	NewPasswordRu              = errors.New("Новый пароль совпадает со старым")
	IncorrectCurrentPasswordRu = errors.New("Неверно указан текущий пароль")

	QuizAddRu = errors.New("Произошла ошибка. Пожалуйста, еще раз проголосуйте")

	NoCommentsRu        = errors.New("У ресторана пока нет отзывов")
	SuccessCleanRu      = errors.New("Корзина очищена")
	NoOrdersRu          = errors.New("Нет заказов")
	NoBasketRu          = errors.New("У Вас нет корзины")
	AlreadyPayedRu      = errors.New("Заказ уже оплачен")
	NoAddFoodToOrderRu  = errors.New("Не удалось добавить блюдо в заказ, попробуйте еще раз")
	FailCleanBasketRu   = errors.New("Не удалось очистить корзину")
	FailCreateCommentRu = errors.New("Не удалось добавить отзыв")
	NoDeleteFoodRu      = errors.New("Не удалось убрать блюдо из заказа, попробуйте еще раз")
	OverDatePromocodeRu = errors.New("Срок действия промокода истек")
	OncePromocodeRu     = errors.New("Данный промокод нельзя использовать повторно")
	SumPromocodeRu      = errors.New("Чтобы применить промокод, сумма должна быть больше")
	NoSetPromocodeRu    = errors.New("Не удалось применить промокод")
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

	OrderAddFood     = errors.New("food was not added to the order")
	OrderSum         = errors.New("sum of the order is null")
	FailCleanBasket  = errors.New("can't clean basket")
	BasketCreate     = errors.New("can't create basket")
	InvalidAddressEn = errors.New("invalid length of address")

	NoComments        = errors.New("restaurant doesn't have comments")
	QuizAdd           = errors.New("answer was not added to the quiz")
	NoBasket          = errors.New("basket doesn't exist")
	AlreadyPayed      = errors.New("order has already been paid")
	OverDatePromocode = errors.New("time of promocode end")
	OncePromocode     = errors.New("this promocode cannot be reused")
	SumPromocode      = errors.New("order sum less than promocode limit")
	NoSetPromocode    = errors.New("promocode could not be applied")

	// Database
	SqlNoRowsUserRelation         = errors.New("no such record exists for \"user\"")
	SqlNoRowsUserRelationAffected = errors.New("no rowы affected for \"user\"")

	SqlNoRowsUnauthAddressRelation = errors.New("no such record exists for unauth_address")
	SqlNoRowsRestaurantRelation    = errors.New("no such record exists for restaurant")
	SqlNoRowsFoodRelation          = errors.New("no such record exists for food")
	SqlNoRowsOrderRelation         = errors.New("no such record exists for \"order\"")
	SqlNoRowsCommentRelation       = errors.New("no such record exists for \"comment\"")
	SqlNoRowsFoodOrderRelation     = errors.New("no such record exists for food_order")
	SqlNoRowsQuizRelation          = errors.New("no such record exists for quiz")
	SqlNoRowsPromocodeRelation     = errors.New("no such record exists for promocode")
	RedisNoData                    = errors.New("no such record exists in Redis")
	NullData                       = errors.New("selected null data from table")
)

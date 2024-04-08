package constants

const (
	SessionCookieName = "session_id"
	CsrfCookieName    = "csrf_token"
)

// Статусы заказа
const (
	Draft    = "draft"
	Payed    = "payed"
	OnTheWay = "on-the-way"
)

const (
	RepoLayer       = "repository"
	UsecaseLayer    = "usecase"
	DeliveryLayer   = "delivery"
	MiddlewareLayer = "middleware"
)

// Название бакетов для minio
const (
	BucketUser = "users"
	BucketRest = "restaurants"
	BucketFood = "foods"
)

// USECASE && REPOSITORY
// Название методов User для логгера
const (
	NameMethodGetUserById           = "GetUserById"
	NameMethodGetUserByEmail        = "GetUserByEmail"
	NameMethodDeleteUserById        = "DeleteUserById"
	NameMethodDeleteUserByEmail     = "DeleteUserByEmail"
	NameMethodCreateUser            = "CreateUser"
	NameMethodUpdateUser            = "UpdateUser"
	NameMethodIsExistById           = "IsExistById"
	NameMethodIsExistByEmail        = "IsExistByEmail"
	NameMethodCheckPassword         = "CheckPassword"
	NameMethodUploadImageByEmail    = "UploadImageByEmail"
	NameMethodGetHashedUserPassword = "GetHashedUserPassword"
	NameMethodUpdateOrder           = "UpdateOrder"
)

const (
	NameMethodGetBasketId = "GetBasketId"
	NameMethodGetBasket   = "GetBasket"
	NameMethodCreateOrder = "CreateOrder"
	NamePayOrder          = "Pay"
)

// Название методов Session для логгера
const (
	NameMethodGetValue  = "GetValue"
	NameMethodSetValue  = "SetValue"
	NameMethodDeleteKey = "DeleteKey"
)

//////////////////////////////////////////////////////

// DELIVERY
// Название методов User для логгера
const (
	NameHandlerUserData    = "UserData"
	NameHandlerUploadImage = "UploadImage"
)

// Название методов auth для логгера
const (
	NameHandlerSignIn  = "SignIn"
	NameHandlerSignUp  = "SignUp"
	NameHandlerSignOut = "SignOut"
)

//////////////////////////////////////////////////////

// MIDDLEWARES
const (
	NameSessionAuthenticationMiddleware = "SessionAuthenticationMiddleware"
	NameCorsMiddleware                  = "CorsMiddleware"
	NameCsrfMiddleware                  = "CorsMiddleware"
)

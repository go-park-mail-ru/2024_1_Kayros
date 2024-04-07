package constants

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
	NameMethodGetBasketId           = "GetBasketId"
	NameMethodGetBasket             = "GetBasket"
	NameMethodCreateOrder           = "CreateOrder"
	NameMethodGetHashedUserPassword = "GetHashedUserPassword"
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
)

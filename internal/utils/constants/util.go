package constants

const (
	SessionCookieName = "session_id"
	CsrfCookieName    = "csrf_token"
)

const (
	ContextCsrf = "need_new_csrf_token"
	Timestamptz = "2006-01-02 15:04:05-07:00"
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

// Название бакетов для minios3
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
	NameMethodAddToOrder            = "AddToOrder"
	NameMethodSetNewPassword        = "SetNewPassword"
)

const (
	NameMethodGetFoodByRest      = "GetFoodByRest"
	NameMethodGetFoodById        = "GetFoodById"
	NameMethodGetBasketId        = "GetBasketId"
	NameMethodGetBasket          = "GetBasket"
	NameMethodCreateOrder        = "CreateOrder"
	NameMethodGetOrders          = "GetOrders"
	NameMethodGetOrderById       = "GetOrderById"
	NameMethodPayOrder           = "Pay"
	NameMethodGetAllRests        = "GetAllRestaurants"
	NameMethodGetRestById        = "GetRestById"
	NameMethodAddFood            = "AddFood"
	NameMethodUpdateCountInOrder = "UpdateCountInOrder"
	NameMethodDeleteFromOrder    = "DeleteFromOrder"
	NameMethodUpdateSum          = "UpdateSum"
	NameMethodGetFoodCount       = "GetFoodCount"
	NameMethodGetFoodPrice       = "GetFoodPrice"
	NameMethodGetOrderSum        = "GetOrderSum"
	NameMethodUpdateStatus       = "UpdateStatus"
	NameMethodUpdateAddress      = "UpdateAddress"
	NameMethodGetFood            = "GetFood"
	NameMethodCleanBasket        = "CleanBasket"
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
	NameHandlerUserData       = "UserData"
	NameHandlerUpdateUser     = "UpdateUser"
	NameHandlerUpdateAddress  = "UpdateAddress"
	NameHandlerUpdatePassword = "UpdatePassword"
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

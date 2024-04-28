package constants

const (
<<<<<<< HEAD
	SessionCookieName     = "session_id"
	CsrfCookieName        = "csrf_token"
	UnauthTokenCookieName = "unauth_token"
=======
	SessionCookieName  = "session_id"
	CsrfCookieName     = "csrf_token"
	UnauthIdCookieName = "unauth_id"

	RequestId   = "request_id"
	XCsrfHeader = "XCSRF_Token"
)

const (
	Timestamptz = "2006-01-02 15:04:05-07:00"
>>>>>>> fix_csrf_test
)

// Статусы заказа
const (
	Draft    = "draft"
	Payed    = "payed"
	OnTheWay = "on-the-way"
)

// Настройка хэширования с помощью Argon2
const (
	HashTime    = 1        // specifies the number of passes over the memory
	HashMemory  = 2 * 1024 // specifies the size of the memory in KiB
	HashThreads = 2
	HashKeylen  = 56
	HashLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

// Название бакетов для minios3
const (
	BucketUser = "users"
	BucketRest = "restaurants"
	BucketFood = "foods"
)

<<<<<<< HEAD
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
	NameMethodDeleteImageByEmail    = "DeleteImageByEmail"
	NameMethodGetHashedUserPassword = "GetHashedUserPassword"
	NameMethodUpdateOrder           = "UpdateOrder"
	NameMethodAddToOrder            = "AddToOrder"
	NameMethodSetNewPassword        = "SetNewPassword"
)
=======
const UploadedFileMaxSize = 10 << 20
>>>>>>> fix_csrf_test

var ValidMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/bmp":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	"image/tiff":    true,
}

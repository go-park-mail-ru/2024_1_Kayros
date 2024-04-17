package constants

const (
	SessionCookieName = "session_id"
	CsrfCookieName    = "csrf_token"
	RequestId         = "request_id"
)

const (
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

const UploadedFileMaxSize = 10 << 20

var ValidMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/bmp":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	"image/tiff":    true,
}

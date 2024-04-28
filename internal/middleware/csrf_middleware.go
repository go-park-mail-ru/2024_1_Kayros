package middleware

import (
	"context"
<<<<<<< HEAD
	"errors"
	"fmt"
=======
	"crypto/sha256"
	"encoding/hex"
>>>>>>> fix_csrf_test
	"net/http"
	"strings"

	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
<<<<<<< HEAD
	"2024_1_kayros/internal/utils/myerrors"
	"2024_1_kayros/internal/utils/regex"
)

var allowedRequestURI = []string{
	"/api/v1/signin", "/api/v1/signup", "/api/v1/user/address", "/api/v1/order/clean", "/api/v1/order/add",
	"/api/v1/order/food/add", "/api/v1/order/food/update_count", "/api/v1/order/food/delete/14", "/api/v1/order/update_address"}
var notAllowedRequestURI = []string{
	"/api/v1/order/pay"}

// CsrfMiddleware проверяет наличие csrf_token в запросе | Метод Signed Double-Submit Cookie
func CsrfMiddleware(handler http.Handler, ucCsrfTokens session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
=======
	"go.uber.org/zap"
)

// Csrf checks for csrf_token availability in the request | Method `Signed Double-Submit Cookie`
func Csrf(handler http.Handler, ucCsrfTokens session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
>>>>>>> fix_csrf_test
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//requestId := functions.GetCtxRequestId(r)
		unauthId := functions.GetCtxUnauthId(r)

<<<<<<< HEAD
		unauthToken := ""
		unauthTokenCookie, err := r.Cookie(cnst.UnauthTokenCookieName)
		if err != nil {
			functions.LogError(logger, requestId, cnst.NameCsrfMiddleware, err, cnst.MiddlewareLayer)
		}
		if unauthTokenCookie != nil {
			unauthToken = unauthTokenCookie.Value
		}
		ctx := context.WithValue(r.Context(), cnst.UnauthTokenCookieName, unauthToken)
		r = r.WithContext(ctx)

		// Будем запрещать доступ к не идемпотентным запросам без валидной сессии
		reqMethod := r.Method
		mutatingMethods := []string{"POST", "PUT", "DELETE"}
		rMethodIsMut := contains(mutatingMethods, reqMethod)
		if rMethodIsMut {
			csrfToken := ""
			csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
			if csrfCookie != nil {
				csrfToken = csrfCookie.Value
			}
			req := regex.RegexURI.MatchString(r.RequestURI)
			fmt.Println(req)
			if errors.Is(err, http.ErrNoCookie) && contains(allowedRequestURI, r.RequestURI) {
				handler.ServeHTTP(w, r)
				return
			} else if err != nil {
				err := errors.New(myerrors.UnauthorizedError)
				functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
				w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusForbidden)
				return
			}

			sessionCookie, err := r.Cookie("session_id")
			sessionId := ""
			if sessionCookie != nil {
				sessionId = sessionCookie.Value
			}
			//if err != nil {
			//	fmt.Println(err, sessionCookie)
			//	err := errors.New(myerrors.UnauthorizedError)
			//	functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
			//	w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusForbidden)
			//	return
			//}

			secretKey := cfg.Server.CsrfSecretKey
			isValid := csrfTokenIsValid(logger, requestId, csrfToken, secretKey, sessionId)
			if !isValid {
				err := errors.New(myerrors.UnauthorizedError)
				functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
				w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusForbidden)
				return
			}
			value, err := ucCsrfTokens.GetValue(r.Context(), alias.SessionKey(csrfToken))
			if err != nil || value == "" {
				err := errors.New(myerrors.UnauthorizedError)
				functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
				w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusForbidden)
				return
			}
		}
=======
		//csrfToken, err := functions.GetCookieCsrfValue(r)
		//if err != nil {
		//	logger.Warn(err.Error(), zap.String(cnst.RequestId, requestId))
		//	w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
		//	return
		//}
		//
		//isValid := csrfTokenIsValid(csrfToken, cfg.Server.CsrfSecretKey)
		//if !isValid {
		//	logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		//	w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
		//	return
		//}
		//
		//xCsrfTokenHeader := r.Header.Get(cnst.XCsrfHeader)
		//if xCsrfTokenHeader != csrfToken {
		//	logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		//	w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
		//	return
		//}
		//
		//_, err = ucCsrfTokens.GetValue(r.Context(), alias.SessionKey(csrfToken))
		//if err != nil {
		//	logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
		//	w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
		//	return
		//}

		ctx := context.WithValue(r.Context(), cnst.UnauthIdCookieName, unauthId)
		r = r.WithContext(ctx)
>>>>>>> fix_csrf_test
		handler.ServeHTTP(w, r)
	})
}

<<<<<<< HEAD
// Функция для проверки наличия элемента в срезе
func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	if regex.RegexURI.MatchString(item) {
		return true
	}
	return false
}

func csrfTokenIsValid(logger *zap.Logger, requestId string, csrfToken string, secretKey string, sessionId string) bool {
	methodName := "csrfTokenIsValid"
	hashData, err := functions.HashCsrf(secretKey, sessionId)
	if err != nil {
		functions.LogError(logger, requestId, methodName, err, cnst.MiddlewareLayer)
		return false
	}
=======
// csrfTokenIsValid | csrf_token consist of 2 parts: hmac and message (hmac it's hash of secretKey and random message).
func csrfTokenIsValid(csrfToken string, secretKey string) bool {
>>>>>>> fix_csrf_test
	parts := strings.Split(csrfToken, ".")
	if len(parts) != 2 {
		return false
	}
	message := parts[1]
	hash := sha256.New()
	_, err := hash.Write([]byte(secretKey + message))
	if err != nil {
		return false
	}
	hmac := hex.EncodeToString(hash.Sum(nil))
	return parts[0] == hmac
}

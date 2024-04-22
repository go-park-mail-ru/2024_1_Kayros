package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
	"2024_1_kayros/internal/utils/myerrors"
	"go.uber.org/zap"
)

var allowedRequestURI = []string{"/api/v1/signin", "/api/v1/signup", "/api/v1/user/address", "/api/v1/order/clean", "/api/v1/order/add",
	"/api/v1/order/food/add", "/api/v1/order/food/update_count", "/order/food/delete/{food_id}"}

// CsrfMiddleware проверяет наличие csrf_token в запросе | Метод Signed Double-Submit Cookie
func CsrfMiddleware(handler http.Handler, ucCsrfTokens session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := ""
		requestIdCtx := r.Context().Value("request_id")
		if requestIdCtx != nil {
			requestId = requestIdCtx.(string)
		}

		// Будем запрещать доступ к не идемпотентным запросам без валидной сессии
		reqMethod := r.Method
		mutatingMethods := []string{"POST", "PUT", "DELETE"}
		rMethodIsMut := contains(mutatingMethods, reqMethod)
		if rMethodIsMut {
			unauthToken := ""
			unauthTokenCookie, err := r.Cookie(cnst.UnauthTokenCookieName)
			if unauthTokenCookie != nil {
				unauthToken = unauthTokenCookie.Value
			}
			ctx := context.WithValue(r.Context(), cnst.UnauthTokenCookieName, unauthToken)
			r = r.WithContext(ctx)

			csrfToken := ""
			csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
			if csrfCookie != nil {
				csrfToken = csrfCookie.Value
			}
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
			if err != nil {
				err := errors.New(myerrors.UnauthorizedError)
				functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
				w = functions.ErrorResponse(w, myerrors.UnauthorizedError, http.StatusForbidden)
				return
			}

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
		handler.ServeHTTP(w, r)
	})
}

// Функция для проверки наличия элемента в срезе
func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
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
	parts := strings.Split(csrfToken, ".")
	if len(parts) != 2 {
		return false
	}
	return hashData == parts[0]
}

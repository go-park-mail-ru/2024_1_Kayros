package middleware

import (
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
		if !rMethodIsMut {
			handler.ServeHTTP(w, r)
			return
		}

		csrfToken := ""
		csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
		if csrfCookie != nil {
			csrfToken = csrfCookie.Value
		}
		if errors.Is(err, http.ErrNoCookie) && (r.RequestURI == "/api/v1/signin" || r.RequestURI == "/api/v1/signup") {
			handler.ServeHTTP(w, r)
			return
		} else if err != nil {
			err := errors.New(myerrors.UnauthorizedError)
			functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
			w = functions.JsonResponse(w, myerrors.UnauthorizedError)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secretKey := cfg.Server.CsrfSecretKey
		isValid := csrfTokenIsValid(csrfToken, secretKey)
		if !isValid {
			err := errors.New(myerrors.UnauthorizedError)
			functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
			w = functions.JsonResponse(w, myerrors.UnauthorizedError)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		value, err := ucCsrfTokens.GetValue(r.Context(), alias.SessionKey(csrfToken))
		if err != nil || value == "" {
			err := errors.New(myerrors.UnauthorizedError)
			functions.LogErrorResponse(logger, requestId, cnst.NameCsrfMiddleware, err, http.StatusForbidden, cnst.MiddlewareLayer)
			w = functions.JsonResponse(w, myerrors.UnauthorizedError)
			w.WriteHeader(http.StatusUnauthorized)
			return
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

func csrfTokenIsValid(csrfToken string, secretKey string) bool {
	parts := strings.Split(csrfToken, ".")
	msg := parts[1]
	newToken := functions.HashData([]byte(secretKey), msg)
	token := string(newToken)
	return token == parts[0]
}
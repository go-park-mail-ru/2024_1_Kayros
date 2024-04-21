package middleware

import (
	"crypto/sha256"
	"encoding/hex"
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

// Csrf checks for csrf_token availability in the request | Method `Signed Double-Submit Cookie`
func Csrf(handler http.Handler, ucCsrfTokens session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId, err := functions.GetCtxRequestId(r)
		if err != nil {
			logger.Error(err.Error())
		}

		// We deny access to non-idempotent requests without a session cookie
		reqMethod := r.Method
		mutatingMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
		rMethodIsMut := contains(mutatingMethods, reqMethod)
		if rMethodIsMut {
			_, err := r.Cookie(cnst.SessionCookieName)
			// We ignore the fact that attacker can sign in/sign out on behalf of the user. It's safe operation.
			if errors.Is(err, http.ErrNoCookie) && (r.RequestURI == "/api/v1/signin" || r.RequestURI == "/api/v1/signup") {
				handler.ServeHTTP(w, r)
				return
			} else if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
				return
			}

			csrfTokenCookieHeader := ""
			csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
			if csrfCookie != nil {
				csrfTokenCookieHeader = csrfCookie.Value
			}
			if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
				return
			}

			isValid := csrfTokenIsValid(csrfTokenCookieHeader, cfg.Server.CsrfSecretKey)
			if !isValid {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
				return
			}

			xCsrfTokenHeader := r.Header.Get(cnst.XCsrfHeader)
			if xCsrfTokenHeader != csrfTokenCookieHeader {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
				return
			}

			_, err = ucCsrfTokens.GetValue(r.Context(), alias.SessionKey(csrfTokenCookieHeader))
			if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
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

func csrfTokenIsValid(csrfToken string, secretKey string) bool {
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

package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"2024_1_kayros/internal/utils/alias"
	"2024_1_kayros/internal/utils/myerrors"
	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

// Csrf checks for csrf_token availability in the request | Method `Signed Double-Submit Cookie`
func Csrf(handler http.Handler, ucSession session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := functions.GetCtxRequestId(r)
		csrfToken, err := functions.GetCookieCsrfValue(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		sessionId, err := functions.GetCookieSessionValue(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}

		isAuthURI := r.RequestURI == "/api/v1/signup" || r.RequestURI == "/api/v1/signin"
		// we don't check csrf token when:
		// 1) user try to authorize
		// 2) user aren't authorized
		if isAuthURI || sessionId == "" {
			handler.ServeHTTP(w, r)
			return
		}

		isValid := csrfTokenIsValid(csrfToken, cfg.Server.CsrfSecretKey)
		if !isValid {
			errMsg := fmt.Sprintf("tokens are not equal %s:%s", invalidCsrfToken, csrfToken)
			logger.Error(errMsg, zap.String(cnst.RequestId, requestId))
			w, err = functions.FlashCookie(r, w, ucSession, &cfg.Redis)
			if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusUnauthorized)
			return
		}

		xCsrfTokenHeader := r.Header.Get(cnst.XCsrfHeader)
		if xCsrfTokenHeader != csrfToken {
			logger.Error(headerCsrfTokenMissing, zap.String(cnst.RequestId, requestId))
			w, err = functions.FlashCookie(r, w, ucSession, &cfg.Redis)
			if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
			return
		}

		_, err = ucSession.GetValue(r.Context(), alias.SessionKey(csrfToken), int32(cfg.DatabaseCsrf))
		if err != nil {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			w, err = functions.FlashCookie(r, w, ucSession, &cfg.Redis)
			if err != nil {
				logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
				w = functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
				return
			}
			w = functions.ErrorResponse(w, myerrors.UnauthorizedRu, http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// csrfTokenIsValid | csrf_token consist of 2 parts: hmac and message (hmac it's hash of secretKey and random message).
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

var (
	invalidCsrfToken       = "csrf-token is invalid"
	headerCsrfTokenMissing = fmt.Sprintf("header '%s' missing", cnst.XCsrfHeader)
)

package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"2024_1_kayros/config"
	"2024_1_kayros/internal/usecase/session"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/functions"
)

// Csrf checks for csrf_token availability in the request | Method `Signed Double-Submit Cookie`
func Csrf(handler http.Handler, ucCsrfTokens session.Usecase, cfg *config.Project, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//requestId := functions.GetCtxRequestId(r)
		unauthId := functions.GetCtxUnauthId(r)
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

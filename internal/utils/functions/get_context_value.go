package functions

import (
	"net/http"

	cnst "2024_1_kayros/internal/utils/constants"
)

func GetCtxRequestId(r *http.Request) string {
	ctxRequestId := r.Context().Value(cnst.RequestId)
	if ctxRequestId == nil {
		return ""
	}
	return ctxRequestId.(string)
}

func GetCtxEmail(r *http.Request) string {
	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		return ""
	}
	return ctxEmail.(string)
}

func GetCtxUnauthId(r *http.Request) string {
	ctxRequestId := r.Context().Value(cnst.UnauthIdCookieName)
	if ctxRequestId == nil {
		return ""
	}
	return ctxRequestId.(string)
}

func GetCookieSessionValue(r *http.Request) (string, error) {
	sessionId := ""
	sessionIdCookie, err := r.Cookie(cnst.SessionCookieName)
	if err != nil {
		return "", err
	}
	sessionId = sessionIdCookie.Value
	return sessionId, nil
}

func GetCookieUnauthIdValue(r *http.Request) (string, error) {
	unauthId := ""
	unauthIdCookie, err := r.Cookie(cnst.UnauthIdCookieName)
	if err != nil {
		return "", err
	}
	unauthId = unauthIdCookie.Value
	return unauthId, nil
}

func GetCookieCsrfValue(r *http.Request) (string, error) {
	csrfToken := ""
	csrfTokenCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		return "", err
	}
	csrfToken = csrfTokenCookie.Value
	return csrfToken, nil
}

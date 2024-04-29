package functions

import (
	"errors"
	"net/http"
	"time"

	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/myerrors"
)

// DeleteCookiesFromDB - method deletes session_id and csrf_token from Redis dbs
func DeleteCookiesFromDB(r *http.Request, ucCsrf session.Usecase, ucSession session.Usecase) error {
	sessionCookie, err := r.Cookie(cnst.SessionCookieName)
	if err != nil {
		return err
	}
	err = ucSession.DeleteKey(r.Context(), alias.SessionKey(sessionCookie.Value))
	if err != nil {
		return err
	}

	csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		return err
	}
	err = ucCsrf.DeleteKey(r.Context(), alias.SessionKey(csrfCookie.Value))
	if err != nil {
		return err
	}
	return nil
}

// CookieExpired - method set cookie expired
func CookieExpired(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, error) {
	sessionCookie, err := r.Cookie(cnst.SessionCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return w, nil
		}
		return w, err
	}
	csrfCookie, err := r.Cookie(cnst.CsrfCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return w, nil
		}
		return w, err
	}

	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionCookie)

	csrfCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, csrfCookie)

	return w, nil
}

func FlashCookie(r *http.Request, w http.ResponseWriter, ucCsrf session.Usecase, ucSession session.Usecase) (http.ResponseWriter, error) {
	err := DeleteCookiesFromDB(r, ucCsrf, ucSession)
	if err != nil {
		if !(errors.Is(err, myerrors.RedisNoData) || errors.Is(err, http.ErrNoCookie)) {
			return w, err
		}
	}

	w, err = CookieExpired(w, r)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			return w, err
		}
	}
	return w, nil
}
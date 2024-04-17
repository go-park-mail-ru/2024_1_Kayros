package functions

import (
	"net/http"

	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
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

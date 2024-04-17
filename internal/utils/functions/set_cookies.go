package functions

import (
	"net/http"
	"time"

	"2024_1_kayros/internal/delivery/auth"
	"2024_1_kayros/internal/usecase/session"
	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"github.com/satori/uuid"
)

const timeExpDur = 14 * 24 * time.Hour

// SetCookie - !!it is necessary to take into account the time zone!!
func SetCookie(w http.ResponseWriter, r *http.Request, ucCsrf session.Usecase, ucSession session.Usecase, email string, secretKey string) (http.ResponseWriter, error) {
	sessionId := uuid.NewV4().String()
	expiration := time.Now().Add(timeExpDur)
	cookie := http.Cookie{
		Name:     cnst.SessionCookieName,
		Value:    sessionId,
		Expires:  expiration,
		HttpOnly: false,
	}
	err := ucSession.SetValue(r.Context(), alias.SessionKey(sessionId), alias.SessionValue(email))
	if err != nil {
		return w, err
	}
	http.SetCookie(w, &cookie)

	csrfToken, err := auth.GenCsrfToken(secretKey, alias.SessionKey(sessionId))
	if err != nil {
		return w, err
	}
	err = ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(email))
	if err != nil {
		return w, err
	}
	csrfCookie := http.Cookie{
		Name:     cnst.CsrfCookieName,
		Value:    csrfToken,
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &csrfCookie)
	return w, nil
}

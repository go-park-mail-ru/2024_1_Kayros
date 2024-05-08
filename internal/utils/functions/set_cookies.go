package functions

import (
	"net/http"
	"time"

	"2024_1_kayros/internal/utils/alias"
	cnst "2024_1_kayros/internal/utils/constants"
	"2024_1_kayros/internal/utils/props"
	"github.com/satori/uuid"
)

const timeExpDur = 14 * 24 * time.Hour

// SetCookie - !!it is necessary to take into account the time zone!!
func SetCookie(w http.ResponseWriter, r *http.Request, props *props.SetCookieProps) (http.ResponseWriter, error) {
	sessionId := uuid.NewV4().String()
	expiration := time.Now().Add(timeExpDur)
	cookie := http.Cookie{
		Name:     cnst.SessionCookieName,
		Value:    sessionId,
		Expires:  expiration,
		HttpOnly: false,
		Path:     "/",
	}
	err := props.UsecaseSession.SetValue(r.Context(), alias.SessionKey(sessionId), alias.SessionValue(props.Email))
	if err != nil {
		return w, err
	}
	http.SetCookie(w, &cookie)

	csrfToken, err := GenerateCsrfToken(props.SecretKey, alias.SessionKey(sessionId))
	if err != nil {
		return w, err
	}
	err = props.UsecaseCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(props.Email))
	if err != nil {
		return w, err
	}
	csrfCookie := http.Cookie{
		Name:     cnst.CsrfCookieName,
		Value:    csrfToken,
		Expires:  expiration,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, &csrfCookie)
	return w, nil
}

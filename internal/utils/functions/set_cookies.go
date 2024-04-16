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

// надо еще учитывать таймзону
func SetCookie(w http.ResponseWriter, r *http.Request, ucCsrf session.Usecase, ucSession session.Usecase, email string, secretKey string) http.ResponseWriter {
	sessionId := uuid.NewV4()
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     cnst.SessionCookieName,
		Value:    sessionId.String(),
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &cookie)

	err := ucSession.SetValue(r.Context(), alias.SessionKey(sessionId.String()), alias.SessionValue(email))
	if err != nil {
		return w
	}

	csrfToken, err := auth.GenCsrfToken(secretKey, alias.SessionKey(sessionId.String()))
	if err != nil {
		return w
	}
	err = ucCsrf.SetValue(r.Context(), alias.SessionKey(csrfToken), alias.SessionValue(email))
	if err != nil {
		return w
	}
	csrfCookie := http.Cookie{
		Name:     cnst.CsrfCookieName,
		Value:    csrfToken,
		Expires:  expiration,
		HttpOnly: false,
	}
	http.SetCookie(w, &csrfCookie)
	return w
}

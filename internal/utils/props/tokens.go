package props

import "2024_1_kayros/internal/usecase/session"

type SetCookieProps struct {
	UsecaseCsrf    session.Usecase
	UsecaseSession session.Usecase
	Email          string
	SecretKey      string
}

func GetSetCookieProps(ucCsrfProps session.Usecase, ucSessionProps session.Usecase, emailProps string, secretKeyProps string) *SetCookieProps {
	return &SetCookieProps{
		UsecaseCsrf:    ucCsrfProps,
		UsecaseSession: ucSessionProps,
		Email:          emailProps,
		SecretKey:      secretKeyProps,
	}
}

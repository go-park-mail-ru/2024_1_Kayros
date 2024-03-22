package dto

// SignInProps структура данных, получаемая с формы авторизации
type SignInProps struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

// SingUpProps структура данных, получаемая с формы регистрации
type SingUpProps struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Phone    string `json:"phone"`
}

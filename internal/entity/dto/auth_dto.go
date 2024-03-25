package dto

// SignInProps структура данных, получаемая с формы авторизации
type SignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

// SignUpProps структура данных, получаемая с формы регистрации
type SignUpDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Phone    string `json:"phone"`
}

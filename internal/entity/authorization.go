package entity

// AuthorizationProps структура данных, получаемая с формы авторизации
type AuthorizationProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegistrationProps структура данных, получаемая с формы регистрации
type RegistrationProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

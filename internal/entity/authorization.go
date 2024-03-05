package entity

type AuthorizationProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

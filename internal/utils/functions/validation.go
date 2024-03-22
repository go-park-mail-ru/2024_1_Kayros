package functions

import "regexp"

// IsValidPassword проверяет пароль на валидность
func IsValidPassword(password string) bool {
	// Проверка на минимальную длину
	if len(password) < 8 {
		return false
	}

	// Проверка на наличие хотя бы одной буквы
	letterRegex := regexp.MustCompile(`[A-Za-z]`)
	if !letterRegex.MatchString(password) {
		return false
	}

	// Проверка на наличие хотя бы одной цифры
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(password) {
		return false
	}

	// Проверка на наличие разрешенных символов
	validCharsRegex := regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	if !validCharsRegex.MatchString(password) {
		return false
	}

	return true
}

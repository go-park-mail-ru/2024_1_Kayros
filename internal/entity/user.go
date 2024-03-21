package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
)

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ImgUrl   string `json:"img_url"`
}

// CheckPassword проверяет пароль, хранящийся в БД с переданным паролем
func (u *User) CheckPassword(password string) bool {
	hashPassword, err := HashData(password)
	if err != nil {
		return false
	}
	return u.Password == hashPassword
}

// SetPassword устанавливает пароль пользователя
func (u *User) SetPassword(password string) (bool, error) {
	newPassword, err := HashData(password)
	if err != nil {
		return false, err
	}
	u.Password = newPassword
	return true, nil
}

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

// HashData хэширует данные с помощью хэш-функции sha256
func HashData(data string) (string, error) {
	hashedPassword := sha256.New()
	_, err := hashedPassword.Write([]byte(data))
	if err != nil {
		return "", errors.New(UnexpectedServerError)
	}
	return hex.EncodeToString(hashedPassword.Sum(nil)), nil
}

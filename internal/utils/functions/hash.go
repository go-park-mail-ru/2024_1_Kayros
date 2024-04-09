package functions

import (
	"golang.org/x/crypto/argon2"
)

// Для секретных данных
func HashData(salt []byte, plainPassword string) []byte {
	hashedPassword := argon2.IDKey([]byte(plainPassword), salt, 1, 2*1024, 2, 56)
	return append(salt, hashedPassword...)
}

func HashCsrf(sessionId string) string {
	return ""
}

package functions

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"2024_1_kayros/internal/utils/myerrors"
	"golang.org/x/crypto/argon2"
)

// HashData Для секретных данных
func HashData(salt []byte, plainPassword string) []byte {
	hashedPassword := argon2.IDKey([]byte(plainPassword), salt, 1, 2*1024, 2, 56)
	return append(salt, hashedPassword...)
}

// HashCsrf хэширует csrf токен
func HashCsrf(secretKey string, sessionId string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(secretKey + sessionId))
	if err != nil {
		return "", errors.New(myerrors.HashedPasswordError)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

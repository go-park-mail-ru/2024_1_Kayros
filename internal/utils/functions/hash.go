package functions

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"2024_1_kayros/internal/utils/myerrors"
)

// HashData хэширует данные с помощью хэш-функции sha256
func HashData(data string) (string, error) {
	hashedPassword := sha256.New()
	_, err := hashedPassword.Write([]byte(data))
	if err != nil {
		return "", errors.New(myerrors.HashedPasswordError)
	}
	return hex.EncodeToString(hashedPassword.Sum(nil)), nil
}
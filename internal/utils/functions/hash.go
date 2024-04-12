package functions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"

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
	randValue, err := generateRandomString(8)
	if err != nil {
		return "", err
	}
	message := sessionId + "!" + randValue
	_, err = hash.Write([]byte(secretKey + message))
	if err != nil {
		return "", errors.New(myerrors.HashedPasswordError)
	}
	csrfToken := hex.EncodeToString(hash.Sum(nil)) + "." + message
	return csrfToken, nil
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

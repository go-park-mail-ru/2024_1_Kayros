package functions

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/argon2"
)

const (
	hashTime    = 1        // specifies the number of passes over the memory
	hashMemory  = 2 * 1024 // specifies the size of the memory in KiB
	hashThreads = 2
	hashKeylen  = 56
	hashLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

// HashData | hashes credentionals using Argon2
func HashData(saltProps []byte, plainPassword string) []byte {
	salt := make([]byte, len(saltProps))
	copy(salt, saltProps)
	hashedPassword := argon2.IDKey([]byte(plainPassword), salt, hashTime, hashMemory, hashThreads, hashKeylen)
	return append(salt, hashedPassword...)
}

// HashCsrf | hashes CSRF-token using sha256
func HashCsrf(secretKey string, sessionId string) (string, error) {
	hash := sha256.New()
	randValue, err := generateRandomString(8)
	if err != nil {
		return "", err
	}
	message := sessionId + "!" + randValue
	_, err = hash.Write([]byte(secretKey + message))
	if err != nil {
		return "", err
	}
	csrfToken := hex.EncodeToString(hash.Sum(nil)) + "." + message
	return csrfToken, nil
}

func generateRandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(hashLetters))))
		if err != nil {
			return "", err
		}
		ret[i] = hashLetters[num.Int64()]
	}
	return string(ret), nil
}

package pkg

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// NUMBERS const numbers
const NUMBERS = "1234567890"

// CHARACTERS const field
const CHARACTERS = "abcdefghijelmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890"

// GenerateRandomString  function
func GenerateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func HashPassword(password string) (string, error) {
	hash, erro := bcrypt.GenerateFromPassword([]byte(password), 0)
	if erro != nil {
		return "", erro
	}
	return string(hash), nil
}

// CompareHash this function compares the hash with the string and returns a boolean value accordingly.
func CompareHash(hash, password string) bool {
	if eror := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); eror != nil {
		return false
	}
	return true
}

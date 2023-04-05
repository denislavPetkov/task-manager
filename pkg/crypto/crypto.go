package crypto

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(input string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(hash)
}

func GetHashedPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hashedPassword), nil
}

func IsHashedPasswordCorrect(password, hashedPassword string) error {
	passwordBytes := []byte(password)
	hasedPasswordBytes := []byte(hashedPassword)

	return bcrypt.CompareHashAndPassword(hasedPasswordBytes, passwordBytes)
}

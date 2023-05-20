package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLength = 32
)

var (
	logger *zap.Logger
)

func init() {
	logger = zap.L().Named("crypto")
}

func GenerateToken() (string, error) {
	buffer := make([]byte, tokenLength)
	_, err := rand.Read(buffer)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate a random token, error: %v", err))
		return "", err
	}

	logger.Info("Generated a random token successfully")

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func GetHashedPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate a hashed password, error: %v", err))
		return "", err
	}

	logger.Info("Generated a hashed password successfully")

	return string(hashedPassword), nil
}

func IsHashedPasswordCorrect(password, hashedPassword string) error {
	passwordBytes := []byte(password)
	hasedPasswordBytes := []byte(hashedPassword)

	return bcrypt.CompareHashAndPassword(hasedPasswordBytes, passwordBytes)
}

func Hash(input string) string {
	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

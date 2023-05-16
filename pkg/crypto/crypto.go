package crypto

import (
	"encoding/base64"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger *zap.Logger
)

func init() {
	logger = zap.L().Named("crypto")
}

func GenerateToken(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate a random token, error: %v", err))
		return "", err
	}

	logger.Info("Generated a random token successfully")

	return base64.StdEncoding.EncodeToString(hash), nil
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

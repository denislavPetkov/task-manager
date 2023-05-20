package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	database "github.com/denislavpetkov/task-manager/pkg/database/user"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	newPasswordHtml = "newPassword.html"
)

func (c *controller) getNewPassword(gc *gin.Context) {
	userEmail := gc.Param("email")
	passwordRecoveryTokenFromUrl := gc.Query(recoveryTokenKey)

	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, newPasswordHtml, gin.H{
		constants.CsrfKey: csrfToken,
		"email":           userEmail,
		"token":           passwordRecoveryTokenFromUrl,
	})
}

func (c *controller) postNewPassword(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	passwordRecoveryTokenFromUrl := gc.PostForm("_token")

	userEmail := gc.PostForm("_email")
	emailBytes, err := base64.StdEncoding.DecodeString(userEmail)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get decode email, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newPasswordHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	emailHash := crypto.Hash(string(emailBytes))
	passwordRecoveryToken, err := c.userDb.Get(context.TODO(), emailHash)
	if err != nil {
		if err.Error() == database.InvalidKeyErr {
			logger.Error(fmt.Sprintf("Password recovery token for %s expired", userEmail))

			gc.Status(http.StatusBadRequest)

			return
		}

		logger.Error(fmt.Sprintf("Failed to check if password recovery token exists in db, error: %v", err))

		gc.Status(http.StatusInternalServerError)

		return
	}

	if passwordRecoveryToken != passwordRecoveryTokenFromUrl {
		logger.Error(fmt.Sprintf("Invalid password recovery token for %s provided", userEmail))

		gc.Status(http.StatusBadRequest)

		return
	}

	userEmail = string(emailBytes)

	password := gc.PostForm("password")
	confirmPassword := gc.PostForm("confirm_password")

	if password != confirmPassword {
		logger.Info("Provided new recovery passwords do not match")

		gc.HTML(http.StatusBadRequest, newPasswordHtml, gin.H{
			errorKey:          "Passwords do not match",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	hashedPassword, err := crypto.GetHashedPassword(password)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get hashed password, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newPasswordHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = c.userDb.Set(context.TODO(), userEmail, hashedPassword, 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update user's password in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newPasswordHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	_, err = c.userDb.Del(context.TODO(), emailHash)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete recovery key from db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newPasswordHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	logger.Info("User created a new password successfully")

	gc.HTML(http.StatusOK, newPasswordHtml, gin.H{successKey: "Created a new password successfully!"})
}

package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"github.com/denislavpetkov/task-manager/pkg/email"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	passwordResetHtml = "recoverPassword.html"
	recoveryTokenKey  = "recoveryToken"
)

func (c *controller) getPasswordRecover(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, passwordResetHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postPasswordRecover(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	userEmail := gc.PostForm("email")

	emailExists, err := c.userDb.Exists(context.TODO(), userEmail)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check if user exists in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, registerHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	if emailExists != 1 {
		logger.Info(fmt.Sprintf("User with %s email does not exist", userEmail))

		gc.HTML(http.StatusOK, passwordResetHtml, gin.H{successKey: "Sent a recovery password email!"})

		return
	}

	emailBase64 := base64.StdEncoding.EncodeToString([]byte(userEmail))

	emailHash := crypto.Hash(userEmail)
	recoveryToken, err := crypto.GenerateToken()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate a token, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = c.userDb.Set(context.TODO(), emailHash, recoveryToken, constants.PasswordRecoveryTokenExpiration)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to set a password recovery token in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = email.SendRecoveryEmail(userEmail, fmt.Sprintf("http://localhost:8081/newPassword/%s?%s=%s", emailBase64, recoveryTokenKey, recoveryToken))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send password recovery email, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	logger.Info(fmt.Sprintf("Sent a password recovery email to user %s", userEmail))

	gc.HTML(http.StatusOK, passwordResetHtml, gin.H{successKey: "Sent a recovery password email!"})
}

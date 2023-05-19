package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/email"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	passwordResetHtml = "recoverPassword.html"
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

	err = email.SendRecoveryEmail(userEmail, fmt.Sprintf("http://localhost:8081/newPassword/%s", emailBase64))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send password recovery email, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	gc.HTML(http.StatusOK, passwordResetHtml, gin.H{successKey: "Sent a recovery password email!"})
}

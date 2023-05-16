package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	changePasswordHtml = "changePassword.html"
	successKey         = "success"
)

func (c *controller) getPasswordChange(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, changePasswordHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postPasswordChange(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	currentPassword := gc.PostForm("currentPassword")
	email := c.getUserFromSession(gc)

	currentHashedPassword, err := c.userDb.Get(context.TODO(), email)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get password from db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, changePasswordHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = crypto.IsHashedPasswordCorrect(currentPassword, currentHashedPassword)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong current password, error: %v", err))

		gc.HTML(http.StatusBadRequest, changePasswordHtml, gin.H{
			errorKey:          "Wrong current password",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	password := gc.PostForm("password")
	confirmPassword := gc.PostForm("confirm_password")

	if password != confirmPassword {
		logger.Info("Provided passwords do not match")

		gc.HTML(http.StatusBadRequest, changePasswordHtml, gin.H{
			errorKey:          "Passwords do not match",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	hashedPassword, err := crypto.GetHashedPassword(password)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get hashed password, error: %v", err))

		gc.HTML(http.StatusInternalServerError, changePasswordHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = c.userDb.Set(context.TODO(), email, hashedPassword, 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to set new password in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, changePasswordHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	logger.Info("Change of password successful")

	gc.HTML(http.StatusCreated, changePasswordHtml, gin.H{successKey: "Change of password successful!"})
}

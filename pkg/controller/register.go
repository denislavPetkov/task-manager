package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"github.com/gin-gonic/gin"

	csrf "github.com/utrack/gin-csrf"
)

const (
	registerHtml = "register.html"
)

func (c *controller) getRegister(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, registerHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postRegister(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	email := gc.PostForm("email")
	if !govalidator.IsEmail(email) {
		logger.Info("Provided invalid email address")

		gc.HTML(http.StatusBadRequest, registerHtml, gin.H{
			errorKey:          "Email not valid",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	password := gc.PostForm("password")
	confirmPassword := gc.PostForm("confirm_password")

	if password != confirmPassword {
		logger.Info("Provided passwords do not match")

		gc.HTML(http.StatusBadRequest, registerHtml, gin.H{
			errorKey:          "Passwords do not match",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	emailExists, err := c.userDb.Exists(context.TODO(), email)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check if user exists in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, registerHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	if emailExists == 1 {
		logger.Info(fmt.Sprintf("User with %s email already exists", email))

		gc.HTML(http.StatusBadRequest, registerHtml, gin.H{
			errorKey:          "Email Address exists",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	hashedPassword, err := crypto.GetHashedPassword(password)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get hashed password, error: %v", err))

		gc.HTML(http.StatusInternalServerError, registerHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = c.userDb.Set(context.TODO(), email, hashedPassword, 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create create a user in db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, registerHtml, gin.H{
			errorKey:          "Server error",
			constants.CsrfKey: csrfToken,
		})

		return
	}

	logger.Info("User registration successful")

	gc.HTML(http.StatusCreated, registerHtml, gin.H{"success": "Registration successful!"})
}

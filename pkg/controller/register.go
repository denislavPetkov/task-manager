package controller

import (
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"github.com/gin-gonic/gin"

	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) getRegister(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, "register.html", gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postRegister(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	email := gc.PostForm("email")
	if !govalidator.IsEmail(email) {
		gc.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error":           "Email not valid",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	password := gc.PostForm("password")
	confirmPassword := gc.PostForm("confirm_password")

	if password != confirmPassword {
		gc.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error":           "Passwords do not match",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	emailExists, err := c.userDb.Exists(context.TODO(), email)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	if emailExists == 1 {
		gc.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error":           "Email Address exists",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	hashedPassword, err := crypto.GetHashedPassword(password)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	err = c.userDb.Set(context.TODO(), email, hashedPassword, 0)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	gc.HTML(http.StatusCreated, "register.html", gin.H{"success": "Registration successful!"})
}

package controller

import (
	"context"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) getPasswordChange(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, "changePassword.html", gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postPasswordChange(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	password := gc.PostForm("password")
	confirmPassword := gc.PostForm("confirm_password")

	if password != confirmPassword {
		gc.HTML(http.StatusBadRequest, "changePassword.html", gin.H{
			"error":           "Passwords do not match",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	hashedPassword, err := crypto.GetHashedPassword(password)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "changePassword.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	email := c.getUserFromSession(gc)

	err = c.userDb.Set(context.TODO(), email, hashedPassword, 0)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "changePassword.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	gc.HTML(http.StatusCreated, "changePassword.html", gin.H{"success": "Change of password successful!"})
}

package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	database "github.com/denislavpetkov/task-manager/pkg/database/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	csrf "github.com/utrack/gin-csrf"
)

const (
	errorKey = "error"

	serverErrorErrMsg      = "Server error"
	wrongCredentialsErrMsg = "Wrong credentials"

	loginHtml = "login.html"
)

func (c *controller) getLogin(gc *gin.Context) {
	session := sessions.Default(gc)
	sessionID := session.Get(constants.SessionIdKey)
	if sessionID != nil {
		gc.Redirect(http.StatusFound, "/tasks")
		gc.Next()
	}

	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, loginHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postLogin(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	email := gc.PostForm("email")
	password := gc.PostForm("password")

	hashedPassword, err := c.userDb.Get(context.TODO(), email)
	if err != nil {
		if err.Error() == database.InvalidKeyErr {
			logger.Error("Invalid login credentials")

			gc.HTML(http.StatusBadRequest, loginHtml, gin.H{
				errorKey:          wrongCredentialsErrMsg,
				constants.CsrfKey: csrfToken,
			})

			return
		}

		logger.Error(fmt.Sprintf("Failed to get password from db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = crypto.IsHashedPasswordCorrect(password, hashedPassword)
	if err != nil {
		logger.Error("Invalid login credentials")

		gc.HTML(http.StatusBadRequest, loginHtml, gin.H{
			errorKey:          wrongCredentialsErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	sessionToken, err := crypto.GenerateToken()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get session token, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	session := sessions.Default(gc)
	session.Set(constants.SessionIdKey, sessionToken)
	session.Set(constants.SessionUserKey, email)

	err = session.Save()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save current session, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	err = c.taskDb.CreateCollection(email)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create a task collection, error: %v", err))

		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			errorKey:          serverErrorErrMsg,
			constants.CsrfKey: csrfToken,
		})

		return
	}

	logger.Info("Login successful")

	gc.SetCookie(constants.CookieUser, email, 0, "", "localhost", true, false)
	gc.Redirect(http.StatusFound, "/tasks")
}

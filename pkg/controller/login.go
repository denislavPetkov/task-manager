package controller

import (
	"context"
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/crypto"
	database "github.com/denislavpetkov/task-manager/pkg/database/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) getLogin(gc *gin.Context) {
	session := sessions.Default(gc)
	sessionID := session.Get(constants.SessionIdKey)
	if sessionID != nil {
		gc.Redirect(http.StatusFound, "/tasks")
	}

	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, "login.html", gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postLogin(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	email := gc.PostForm("email")
	password := gc.PostForm("password")

	hashedPassword, err := c.userDb.Get(context.TODO(), email)
	if err != nil {
		if err.Error() == database.InvalidKeyErr {
			gc.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error":           "Wrong credentials",
				constants.CsrfKey: csrfToken,
			})
			return
		}
		gc.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error":           "Server error",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	err = crypto.IsHashedPasswordCorrect(password, hashedPassword)
	if err != nil {
		gc.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error":           "Wrong credentials",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	sessionToken := crypto.GenerateToken(email)

	session := sessions.Default(gc)
	session.Set(constants.SessionIdKey, sessionToken)
	session.Set(constants.SessionUserKey, email)
	session.Save()

	err = c.taskDb.CreateCollection(email)
	if err != nil {
		gc.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error":           "Wrong credentials",
			constants.CsrfKey: csrfToken,
		})
		return
	}

	gc.Redirect(http.StatusFound, "/tasks")
}

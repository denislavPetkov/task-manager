package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (c *controller) getLogout(gc *gin.Context) {
	session := sessions.Default(gc)
	session.Options(sessions.Options{
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	session.Clear()
	session.Save()

	gc.Redirect(http.StatusFound, "/login")
}

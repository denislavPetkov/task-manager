package controller

import (
	"fmt"
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

	err := session.Save()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save current session, error: %v", err))
		return
	}

	logger.Info("Logout successful")

	gc.Redirect(http.StatusFound, "/login")
}

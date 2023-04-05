package middleware

import (
	"net/http"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get(constants.SessionIdKey)
		if sessionID == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		sessionUser := session.Get(constants.SessionUserKey)
		if sessionUser == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

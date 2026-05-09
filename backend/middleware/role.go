package middleware

import (
	"net/http"
	"ocs-room-booking/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			utils.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}

		utils.Error(c, http.StatusForbidden, "Forbidden: insufficient permissions")
		c.Abort()
	}
}

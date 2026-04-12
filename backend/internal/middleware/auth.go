package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"taskflow/internal/utils"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, "Missing Authorization header", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendError(c, http.StatusUnauthorized, "Invalid Authorization header format", nil)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		// Set the user ID in the context for handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

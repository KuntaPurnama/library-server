package middleware

import (
	"github.com/gin-gonic/gin"
	"library-backend/internal/helpers"
	"net/http"
)

func ValidateTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		token := authHeader[len("Bearer "):]
		isValid, _ := helpers.ValidateToken(token)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

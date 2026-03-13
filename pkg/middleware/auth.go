package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InternalAuth validates the X-Internal-API-Key header.
func InternalAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Internal-API-Key")

		if expectedKey == "" || apiKey != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

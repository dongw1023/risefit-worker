package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InternalAuth validates the X-Internal-API-Key header or ensures the request is from Cloud Tasks.
func InternalAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Internal-API-Key")
		// Google Cloud Tasks adds this header and it's stripped for external requests on Cloud Run.
		isCloudTask := c.GetHeader("X-CloudTasks-QueueName") != ""

		if !isCloudTask && (expectedKey == "" || apiKey != expectedKey) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

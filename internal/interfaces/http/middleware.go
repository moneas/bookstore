package http

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moneas/bookstore/internal/application"
	"github.com/moneas/bookstore/internal/util"
)

func BasicAuthMiddleware(userService *application.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		payload, err := util.DecodeBase64(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
			c.Abort()
			return
		}

		creds := strings.SplitN(payload, ":", 2)
		if len(creds) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		user, err := userService.Authenticate(creds[0], creds[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func Authenticate(userService *application.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Split "Basic <base64-encoded-credentials>"
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Decode the base64 credentials
		credentials, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 credentials"})
			c.Abort()
			return
		}

		// Split the credentials into username and password
		credParts := strings.SplitN(string(credentials), ":", 2)
		if len(credParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization value"})
			c.Abort()
			return
		}
		username, password := credParts[0], credParts[1]

		// Verify the username and password
		user, err := userService.Authenticate(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		// Store the user ID in the context
		c.Set("user_id", user.ID)
		c.Next()
	}
}

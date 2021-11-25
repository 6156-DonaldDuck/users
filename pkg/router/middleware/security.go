package middleware

import (
	"fmt"
	"net/http"

	"github.com/6156-DonaldDuck/users/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Response struct {
	message string
}

// New security middleware
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		whiteListMap := map[string][]string{
			"/swagger/*any":                 {"GET"},
			"/api/v1/login/google/url":      {"GET"},
			"/api/v1/login/google/callback": {"POST"},
			"/api/v1/users":                 {"GET", "POST"},
			"/api/v1/compositions":          {"GET", "POST"},
			"/api/v1/users/:userId":         {"GET", "PUT"},
		}

		// Figure out whether current request is in white list.
		allowed := false
		allowedMethods, ok := whiteListMap[c.FullPath()]
		if ok {
			for _, method := range allowedMethods {
				if method == c.Request.Method {
					allowed = true
					break
				}
			}
		}
		// if request is in whitelist
		if allowed {
			c.Next()
		} else {
			// Get access token
			accessToken := c.GetHeader("Authorization")

			// if token is empty
			if accessToken == "" {
				err := fmt.Errorf("access token should not be empty")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}
			// if token does not exist
			token := auth.TokenStoreInstance.GetToken(accessToken)
			if token == nil {
				err := fmt.Errorf("token not found for access token=%s", accessToken)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}
			// Allowed
			c.Next()
		}
	}
}

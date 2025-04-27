package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for authentication cookie
		cookie, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the cookie (this is just a placeholder, implement your own logic)
		if cookie == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("Cookie:", cookie)

		userID, err := repository.GetUserByID(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// If the cookie is valid, proceed to the next handler
		c.Set("user_id", userID)
		c.Next()
	}
}
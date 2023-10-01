package funcs

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")

		if tokenHeader == "" {
			c.JSON(401, gin.H{"status": "Unauthorized", "error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := parseToken(tokenHeader)
		if err != nil {
			c.JSON(401, gin.H{"status": "Unauthorized", "error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

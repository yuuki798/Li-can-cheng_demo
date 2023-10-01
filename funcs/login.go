package funcs

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		var user User

		// 从请求中解析用户数据
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid JSON format"})
			return
		}

		// 查询数据库中的密码
		hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)

		if err != nil {
			c.JSON(404, gin.H{"status": "NotFound", "error": "Username not found or password mismatch"})
			return
		}

		// 验证密码
		// ... 在登录函数中 ...
		if hashedPassword == fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password))) {
			token, err := generateToken(user.Username)
			if err != nil {
				c.JSON(500, gin.H{"status": "InternalServerError", "error": "Could not generate token"})
				return
			}
			c.JSON(200, gin.H{"status": "OK", "token": token})
			return
		}

		// 为简化起见，这里没有使用JWT或其他身份验证机制
		c.JSON(200, gin.H{"status": "OK", "message": "Successfully logged in!"})
	})
}

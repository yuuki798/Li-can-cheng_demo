package funcs

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, db *sql.DB) {
	r.POST("/register", func(c *gin.Context) {
		var user User

		// 从请求中解析用户数据
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid JSON format"})
			return
		}

		// 密码加密（这只是一个简单的例子，请在实际中使用更安全的方法如bcrypt）
		hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))

		// 将用户信息存入数据库
		_, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", user.Username, hashedPassword)
		if err != nil {
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "OK", "message": "Successfully registered!"})
	})
}

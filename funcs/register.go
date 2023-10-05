package funcs

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Register(r *gin.Engine, db *sql.DB) {
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		// 从请求中解析用户数据
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, RegisterResponse{Status: "BadRequest", Error: "Invalid JSON format"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, RegisterResponse{Status: "InternalServerError", Error: "Could not hash password"})
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", req.Username, string(hashedPassword))
		if err != nil {
			c.JSON(500, RegisterResponse{Status: "InternalServerError", Error: err.Error()})
			return
		}
		c.JSON(200, RegisterResponse{Status: "OK"})
	})
}
